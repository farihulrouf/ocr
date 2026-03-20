package main

import (
	"context"
	"log"
	"ocr-saas-backend/configs"
	"ocr-saas-backend/internal/service/ocr"

	"github.com/google/uuid"
)

var aiLimiter = make(chan struct{}, 3) // max 2 AI request bersamaan
func main() {
	configs.LoadConfig()
	configs.InitS3()

	//configs.InitMinioConfig() // <-- pastikan ditambahkan
	configs.ConnectDB()
	//configs.SeedDatabase(configs.DB)
	configs.ConnectRedis()

	ctx := context.Background()
	log.Println("OCR Worker started...")

	// 🔥 1️⃣ RECOVER JOB DULU
	recoverStuckJobs()

	for {
		receiptID, err := configs.RedisClient.
			BRPopLPush(ctx, "ocr:queue", "ocr:processing", 0).
			Result()

		if err != nil {
			log.Println("[ERROR] Redis BRPopLPush:", err)
			continue
		}

		log.Println("[DEBUG] Got receiptID:", receiptID)

		// Tandai PROCESSING & set started_at
		id, _ := uuid.Parse(receiptID)
		if err := ocr.SetOCRJobStatus(id, "PROCESSING", ""); err != nil {
			log.Println("[ERROR] Failed to set PROCESSING:", err)
			continue
		}

		aiLimiter <- struct{}{} // lock AI
		err = ocr.ProcessOCRString(receiptID)
		<-aiLimiter // unlock

		if err != nil {
			// FAILED & set finished_at
			ocr.SetOCRJobStatus(id, "FAILED", err.Error())
			ocr.MarkAsFailed(receiptID, err.Error())
			log.Println("[ERROR] OCR failed:", err)
		} else {
			// DONE & set finished_at
			ocr.SetOCRJobStatus(id, "DONE", "")
			ocr.MarkAsSuccess(receiptID)
			log.Println("[DEBUG] OCR processed:", receiptID)
		}

		// Hapus dari processing list
		configs.RedisClient.LRem(ctx, "ocr:processing", 1, receiptID)
	}
}

func recoverStuckJobs() {
	ctx := context.Background()
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
