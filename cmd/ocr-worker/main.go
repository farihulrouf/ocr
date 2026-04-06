package main

import (
	"context"
	"log"
	"ocr-saas-backend/configs"
	"ocr-saas-backend/internal/service/ocr"
	"time"

	"github.com/google/uuid"
)

var aiLimiter = make(chan struct{}, 3)

func main() {
	cfg := configs.LoadConfig()

	configs.ConnectDB(cfg)
	configs.ConnectRedis(cfg)
	configs.InitS3(cfg)

	ctx := context.Background()
	log.Println("👷 OCR Worker is ready and listening to 'ocr:queue'...")

	for {
		// BRPopLPush: Ambil dari queue, pindah ke processing (Reliable Queue)
		receiptID, err := configs.RedisClient.
			BRPopLPush(ctx, "ocr:queue", "ocr:processing", 0).
			Result()

		if err != nil {
			log.Println("❌ Redis Queue Error:", err)
			time.Sleep(2 * time.Second)
			continue
		}

		log.Println("📝 Processing Receipt ID:", receiptID)
		id, _ := uuid.Parse(receiptID)

		// Set status PROCESSING ke DB
		ocr.SetOCRJobStatus(id, "PROCESSING", "")

		// AI Concurrency Limit
		aiLimiter <- struct{}{}
		err = ocr.ProcessOCRString(receiptID)
		<-aiLimiter

		if err != nil {
			handleFailure(ctx, receiptID, id, err)
		} else {
			handleSuccess(ctx, receiptID, id)
		}
	}
}

func handleSuccess(ctx context.Context, receiptID string, id uuid.UUID) {
	ocr.SetOCRJobStatus(id, "DONE", "")
	ocr.MarkAsSuccess(receiptID)
	configs.RedisClient.LRem(ctx, "ocr:processing", 1, receiptID)
	log.Println("✅ OCR Success:", receiptID)
}

func handleFailure(ctx context.Context, receiptID string, id uuid.UUID, err error) {
	log.Println("❌ OCR Failed:", err)
	// Logika retry atau pindah ke dead letter queue bisa di sini
	ocr.SetOCRJobStatus(id, "FAILED", err.Error())
	configs.RedisClient.LRem(ctx, "ocr:processing", 1, receiptID)
}
