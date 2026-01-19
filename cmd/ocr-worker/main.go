package main

import (
	"context"
	"log"
	"ocr-saas-backend/configs"
	"ocr-saas-backend/internal/service/ocr"
)

func main() {
	configs.LoadConfig()
	configs.ConnectDB()
	configs.ConnectRedis()

	ctx := context.Background()
	log.Println("OCR Worker started...")

	for {
		result, err := configs.RedisClient.BLPop(ctx, 0, "ocr:queue").Result()
		if err != nil {
			log.Println("[ERROR] Redis BLPop:", err)
			continue
		}

		receiptID := result[1]
		log.Println("[DEBUG] Got receiptID from queue:", receiptID)

		if err := ocr.ProcessOCRString(receiptID); err != nil {
			log.Println("[ERROR] OCR failed for", receiptID, ":", err)
		} else {
			log.Println("[DEBUG] OCR processed for", receiptID)
		}
	}
}
