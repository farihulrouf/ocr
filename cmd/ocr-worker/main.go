package main

import (
	"context"
	"log"
	"ocr-saas-backend/configs"
	"ocr-saas-backend/internal/service/ocr"
	"time"

	"github.com/google/uuid"
)

var aiLimiter = make(chan struct{}, 3) // max 3 AI request bersamaan

func main() {
	// ✅ Load config
	cfg := configs.LoadConfig()

	// ✅ Init semua dependency pakai cfg
	configs.InitS3(cfg)
	configs.ConnectDB(cfg)
	configs.ConnectRedis(cfg)

	ctx := context.Background()
	log.Println("🚀 OCR Worker started...")

	// 🔥 Recover stuck jobs dulu
	recoverStuckJobs(ctx)

	for {
		receiptID, err := configs.RedisClient.
			BRPopLPush(ctx, "ocr:queue", "ocr:processing", 0).
			Result()

		if err != nil {
			log.Println("[ERROR] Redis BRPopLPush:", err)
			continue
		}

		log.Println("[DEBUG] Got receiptID:", receiptID)

		id, _ := uuid.Parse(receiptID)

		// ✅ Set status PROCESSING
		if err := ocr.SetOCRJobStatus(id, "PROCESSING", ""); err != nil {
			log.Println("[ERROR] Failed set PROCESSING:", err)
			continue
		}

		// 🔁 Retry count
		retryKey := "ocr:retry:" + receiptID
		retryCount, _ := configs.RedisClient.Get(ctx, retryKey).Int()

		// 🔥 Limit AI concurrency
		aiLimiter <- struct{}{}
		err = ocr.ProcessOCRString(receiptID)
		<-aiLimiter

		if err != nil {
			log.Println("[ERROR] OCR failed:", err)

			if retryCount < 3 {
				configs.RedisClient.Incr(ctx, retryKey)
				configs.RedisClient.Expire(ctx, retryKey, time.Hour)

				log.Println("[RETRY] Requeue:", receiptID, "retry:", retryCount+1)

				// balik ke queue
				configs.RedisClient.LPush(ctx, "ocr:queue", receiptID)

				ocr.SetOCRJobStatus(id, "PROCESSING", "retrying...")
			} else {
				log.Println("[DEAD] Move to dead queue:", receiptID)

				configs.RedisClient.LPush(ctx, "ocr:dead", receiptID)

				ocr.SetOCRJobStatus(id, "FAILED", err.Error())
				ocr.MarkAsFailed(receiptID, err.Error())
			}

		} else {
			// ✅ Success
			ocr.SetOCRJobStatus(id, "DONE", "")
			ocr.MarkAsSuccess(receiptID)

			log.Println("[SUCCESS] OCR processed:", receiptID)
		}

		// ✅ Remove dari processing list
		configs.RedisClient.LRem(ctx, "ocr:processing", 1, receiptID)
	}
}

func recoverStuckJobs(ctx context.Context) {
	jobs, err := configs.RedisClient.LRange(ctx, "ocr:processing", 0, -1).Result()
	if err != nil {
		log.Println("[ERROR] Recovery failed:", err)
		return
	}

	for _, job := range jobs {
		log.Println("[RECOVERY] Re-queue:", job)

		configs.RedisClient.LRem(ctx, "ocr:processing", 1, job)
		configs.RedisClient.LPush(ctx, "ocr:queue", job)
	}
}
