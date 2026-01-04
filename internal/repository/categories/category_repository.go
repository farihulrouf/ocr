package categories

import (
	"ocr-saas-backend/configs"
	"ocr-saas-backend/internal/models"

	"github.com/google/uuid"
)

func ListCategories(
	tenantID uuid.UUID,
	page, pageSize int,
) ([]models.AccountCategory, int64, error) {

	var rows []models.AccountCategory
	var total int64

	db := configs.DB.
		Model(&models.AccountCategory{}).
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

func GetCategoryByID(
	tenantID, id uuid.UUID,
) (*models.AccountCategory, error) {

	var cat models.AccountCategory

	err := configs.DB.
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&cat).Error

	if err != nil {
		return nil, err
	}

	return &cat, nil
}

func CreateCategory(cat *models.AccountCategory) error {
	return configs.DB.Create(cat).Error
}

func UpdateCategory(cat *models.AccountCategory) error {
	return configs.DB.Save(cat).Error
}

func DeleteCategory(cat *models.AccountCategory) error {
	return configs.DB.Delete(cat).Error // soft delete
}
