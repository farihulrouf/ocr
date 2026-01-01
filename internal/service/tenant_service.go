package service

import (
	"fmt"
	"ocr-saas-backend/internal/models"
	"ocr-saas-backend/internal/repository"
)

func GetTenantInfo(tenantID string) (*models.Tenant, error) {
	return repository.GetTenantByID(tenantID)
}

func UpdateTenantInfo(tenantID string, data map[string]interface{}) error {
	return repository.UpdateTenantInfo(tenantID, data)
}

func GetTenantSettings(tenantID string) (*models.CompanySetting, error) {
	return repository.GetTenantSettings(tenantID)
}
func GetTenantSubscription(tenantID string) (*models.Tenant, error) {
	return repository.GetTenantSubscription(tenantID)
}

func CreateUpgradeCheckoutURL(tenantID string, planID string) (string, error) {
	// cek apakah plan valid
	_, err := repository.GetPlanByID(planID)
	if err != nil {
		return "", err
	}

	// generate dummy checkout URL
	checkoutURL := fmt.Sprintf(
		"https://pay.ocr-saas.com/checkout?tenant=%s&plan=%s",
		tenantID, planID,
	)

	return checkoutURL, nil
}

func GetAllTenants(page, pageSize int, q, sort string) ([]models.Tenant, int64, error) {
	return repository.GetAllTenants(page, pageSize, q, sort)
}
