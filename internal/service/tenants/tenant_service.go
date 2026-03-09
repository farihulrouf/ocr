package tenants

import (
	"fmt"
	"ocr-saas-backend/internal/models"
	"ocr-saas-backend/internal/repository/tenants"
)

func GetTenantInfo(tenantID string) (*models.Tenant, error) {
	return tenants.GetTenantByID(tenantID)
}

func UpdateTenantInfo(tenantID string, data map[string]interface{}) error {
	return tenants.UpdateTenantInfo(tenantID, data)
}

func GetTenantSettings(tenantID string) (*models.CompanySetting, error) {
	return tenants.GetTenantSettings(tenantID)
}
func GetTenantSubscription(tenantID string) (*models.Tenant, error) {
	return tenants.GetTenantSubscription(tenantID)
}

func CreateUpgradeCheckoutURL(tenantID string, planID string) (string, error) {
	// cek apakah plan valid
	_, err := tenants.GetPlanByID(planID)
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
	return tenants.GetAllTenants(page, pageSize, q, sort)
}
