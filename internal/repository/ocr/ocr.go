package ocr

import (
	"ocr-saas-backend/configs"
	"ocr-saas-backend/internal/models"

	"github.com/google/uuid"
)

func CreateReceipt(r *models.Receipt) error {
	return configs.DB.Create(r).Error
}

func GetReceiptByID(id uuid.UUID) (*models.Receipt, error) {
	var r models.Receipt
	if err := configs.DB.First(&r, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

func UpdateReceipt(r *models.Receipt) error {
	return configs.DB.Save(r).Error
}
