package service

import (
	"ocr-saas-backend/internal/repository"

	"github.com/google/uuid"
)

func GetAllUsers(
	tenantID uuid.UUID, page, pageSize int, q, sort string) (interface{}, error) {
	data, total, err := repository.GetAllUsers(tenantID, page, pageSize, q, sort)
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

func GetDetail(tenantID, userID uuid.UUID) (map[string]interface{}, error) {
	return repository.GetUserDetail(tenantID, userID)
}

func UpdateUser(tenantID, userID uuid.UUID, role string, deptID *uuid.UUID) (map[string]interface{}, error) {

	err := repository.UpdateUser(tenantID, userID, role, deptID)
	if err != nil {
		return nil, err
	}

	// konsisten dengan gaya kamu → return response simple
	return map[string]interface{}{
		"message": "User updated",
		"status":  "success",
	}, nil
}
