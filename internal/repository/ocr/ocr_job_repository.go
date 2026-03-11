package ocr

import (
	"ocr-saas-backend/configs"
	"ocr-saas-backend/internal/models"
	"time"

	"github.com/google/uuid"
)

func CreateOCRJob(job *models.OCRJob) error {
	return configs.DB.Create(job).Error
}

func GetOCRJobByID(id uuid.UUID) (*models.OCRJob, error) {
	var job models.OCRJob
	if err := configs.DB.First(&job, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

func GetOCRJobByReceiptID(receiptID uuid.UUID) (*models.OCRJob, error) {
	var job models.OCRJob
	if err := configs.DB.First(&job, "receipt_id = ?", receiptID).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

func UpdateOCRJob(job *models.OCRJob) error {
	return configs.DB.Save(job).Error
}

func UpdateOCRJobStatus(receiptID uuid.UUID, status string) error {
	update := map[string]interface{}{
		"status": status,
	}
	if status == "PROCESSING" {
		update["started_at"] = time.Now()
	}
	if status == "DONE" || status == "FAILED" {
		update["finished_at"] = time.Now()
	}
	return configs.DB.
		Model(&models.OCRJob{}).
		Where("receipt_id = ?", receiptID).
		Updates(update).Error
}

func UpdateOCRJobFailed(receiptID uuid.UUID, errMsg string) error {
	return configs.DB.
		Model(&models.OCRJob{}).
		Where("receipt_id = ?", receiptID).
		Updates(map[string]interface{}{
			"status":        "FAILED",
			"error_message": errMsg,
			"finished_at":   time.Now(),
		}).Error
}
