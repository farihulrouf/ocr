package payments

import (
	"ocr-saas-backend/configs"
	"ocr-saas-backend/internal/models"

	"github.com/google/uuid"
)

func ListPayments(
	tenantID uuid.UUID,
	page, pageSize int,
) ([]models.PaymentMethod, int64, error) {

	var rows []models.PaymentMethod
	var total int64

	db := configs.DB.
		Model(&models.PaymentMethod{}).
		Where("tenant_id = ?", tenantID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	err := db.
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&rows).Error

	return rows, total, err
}

func GetPaymentByID(
	tenantID, id uuid.UUID,
) (*models.PaymentMethod, error) {

	var cat models.PaymentMethod

	err := configs.DB.
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&cat).Error

	if err != nil {
		return nil, err
	}

	return &cat, nil
}

func CreatePayments(cat *models.PaymentMethod) error {
	return configs.DB.Create(cat).Error
}

func UpdatePayments(cat *models.PaymentMethod) error {
	return configs.DB.Save(cat).Error
}

func DeletePayments(
	tenantID, id uuid.UUID,
) error {
	return configs.DB.
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&models.PaymentMethod{}).Error
}
