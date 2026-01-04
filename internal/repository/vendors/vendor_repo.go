package vendors

import (
	"ocr-saas-backend/configs"
	"ocr-saas-backend/internal/models"

	"github.com/google/uuid"
)

func ListVendor(
	tenantID uuid.UUID,
	page, pageSize int,
) ([]models.VendorMaster, int64, error) {

	var rows []models.VendorMaster
	var total int64

	db := configs.DB.
		Model(&models.VendorMaster{}).
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

func GetVendorByID(
	tenantID, id uuid.UUID,
) (*models.VendorMaster, error) {

	var cat models.VendorMaster

	err := configs.DB.
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&cat).Error

	if err != nil {
		return nil, err
	}

	return &cat, nil
}

func CreateVendor(cat *models.VendorMaster) error {
	return configs.DB.Create(cat).Error
}

func UpdateVendor(cat *models.VendorMaster) error {
	return configs.DB.Save(cat).Error
}

func DeleteVendor(
	tenantID, id uuid.UUID,
) error {
	return configs.DB.
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&models.VendorMaster{}).Error
}
