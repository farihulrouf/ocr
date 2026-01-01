package service

import (
	"ocr-saas-backend/internal/models"
	"ocr-saas-backend/internal/repository"

	"github.com/google/uuid"
)

func GetAllDepartments(page, pageSize int, q, sort string) (interface{}, error) {
	data, total, err := repository.GetAllDepartments(page, pageSize, q, sort)
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"data": data,
		"meta": map[string]interface{}{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
		"status": "success",
	}

	return response, nil
}

func CreateDepartment(tenantID uuid.UUID, name string) (*models.Department, error) {
	dept := &models.Department{
		TenantID: tenantID,
		Name:     name,
		Code:     "", // optional, bisa auto generate nanti
	}

	err := repository.CreateDepartment(dept)
	if err != nil {
		return nil, err
	}

	return dept, nil
}
