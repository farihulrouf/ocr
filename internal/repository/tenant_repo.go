package repository

import (
	"ocr-saas-backend/configs"
	"ocr-saas-backend/internal/models"
)

// Ambil tenant berdasarkan ID
func GetTenantByID(tenantID string) (*models.Tenant, error) {
	var tenant models.Tenant

	err := configs.DB.
		Where("id = ?", tenantID).
		First(&tenant).Error

	return &tenant, err
}
