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

func UpdateTenantInfo(tenantID string, data map[string]interface{}) error {
	return configs.DB.Model(&models.Tenant{}).
		Where("id = ?", tenantID).
		Updates(data).Error
}

func GetTenantSettings(tenantID string) (*models.CompanySetting, error) {
	var settings models.CompanySetting
	err := configs.DB.
		Where("tenant_id = ?", tenantID).
		First(&settings).Error

	return &settings, err
}

func GetTenantSubscription(tenantID string) (*models.Tenant, error) {
	var tenant models.Tenant

	err := configs.DB.
		Preload("SubscriptionPlan").
		Where("id = ?", tenantID).
		First(&tenant).Error

	return &tenant, err
}

func GetPlanByID(planID string) (*models.SubscriptionPlan, error) {
	var plan models.SubscriptionPlan
	err := configs.DB.Where("id = ?", planID).First(&plan).Error
	return &plan, err
}

func CountReceiptsByTenant(tenantID string) (int64, error) {
	var count int64
	err := configs.DB.Model(&models.Receipt{}).
		Where("tenant_id = ?", tenantID).
		Count(&count).Error
	return count, err
}

func CountUsersByTenant(tenantID string) (int64, error) {
	var count int64
	err := configs.DB.Model(&models.User{}).
		Where("tenant_id = ?", tenantID).
		Count(&count).Error
	return count, err
}
