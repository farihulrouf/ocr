package tax

import (
	"ocr-saas-backend/configs"
	"ocr-saas-backend/internal/models"

	"github.com/google/uuid"
)

func GetTaxRates(
	tenantID uuid.UUID,
	page, pageSize int,
) ([]models.TaxRate, int64, error) {

	var rows []models.TaxRate
	var total int64

	db := configs.DB.
		Model(&models.TaxRate{}).
		Where("tenant_id = ?", tenantID)

	// count
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

func CreateTaxRate(rate *models.TaxRate) error {
	return configs.DB.Create(rate).Error
}

func GetTaxRateByID(
	tenantID, id uuid.UUID,
) (*models.TaxRate, error) {

	var rate models.TaxRate

	err := configs.DB.
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&rate).Error

	if err != nil {
		return nil, err
	}

	return &rate, nil
}

func UpdateTaxRate(rate *models.TaxRate) error {
	return configs.DB.Save(rate).Error
}

func DeleteTaxRate(
	tenantID, id uuid.UUID,
) error {
	return configs.DB.
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&models.TaxRate{}).Error
}
