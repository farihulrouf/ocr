package ocr

import "ocr-saas-backend/internal/redis"

const OCRQueue = "ocr:queue"

func EnqueueOCR(receiptID string) error {
	return redis.Client.
		RPush(redis.Ctx, OCRQueue, receiptID).
		Err()
}
