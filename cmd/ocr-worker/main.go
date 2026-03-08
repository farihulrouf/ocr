package main

import (
	"context"
	"log"
	"ocr-saas-backend/configs"
	"ocr-saas-backend/internal/service/ocr"
)

func main() {
	configs.LoadConfig()
	configs.InitS3()

	//configs.InitMinioConfig() // <-- pastikan ditambahkan
	configs.ConnectDB()
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

		// 🔵 tandai mulai
		if err := ocr.MarkAsProcessing(receiptID); err != nil {
			log.Println("[ERROR] Failed to mark processing:", err)
			continue
		}

		if err := ocr.ProcessOCRString(receiptID); err != nil {

			// 🔴 kalau gagal
			ocr.MarkAsFailed(receiptID, err.Error())
			log.Println("[ERROR] OCR failed:", err)

		} else {

			// 🟢 kalau sukses
			ocr.MarkAsSuccess(receiptID)
			log.Println("[DEBUG] OCR processed:", receiptID)
		}

		// ✅ HAPUS dari processing list setelah selesai
		configs.RedisClient.LRem(ctx, "ocr:processing", 1, receiptID)
	}
}

func recoverStuckJobs() {
	ctx := context.Background()

	jobs, err := configs.RedisClient.
		LRange(ctx, "ocr:processing", 0, -1).
		Result()

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
