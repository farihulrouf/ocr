package service

import (
	"ocr-saas-backend/internal/models"
	"ocr-saas-backend/internal/repository"
)

func GetTenantInfo(tenantID string) (*models.Tenant, error) {
	return repository.GetTenantByID(tenantID)
}
