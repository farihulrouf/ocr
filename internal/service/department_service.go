package service

import (
	"errors"
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

func GetDepartmentDetail(id uuid.UUID) (interface{}, error) {
	dept, err := repository.GetDepartmentByID(id)
	if err != nil {
		return nil, err
	}

	users, err := repository.GetUsersByDepartment(id)
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"status": "success",
		"id":     dept.ID,
		"name":   dept.Name,
		"users":  users,
	}

	return response, nil
}

func UpdateDepartment(id string, name string) error {
	deptID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid UUID")
	}

	// cek apakah department ada
	dept, err := repository.GetDepartmentByID(deptID)
	if err != nil || dept == nil {
		return errors.New("department not found")
	}

	// update
	if err := repository.UpdateDepartmentByID(deptID, name); err != nil {
		return errors.New("failed to update department")
	}

	return nil
}

func DeleteDepartment(id uuid.UUID) error {

	// Check if department exists
	_, err := repository.GetDepartmentByID(id)
	if err != nil {
		return errors.New("department not found")
	}

	// Check if department still has users
	users, err := repository.GetUsersByDepartment(id)
	if err != nil {
		return err
	}

	if len(users) > 0 {
		return errors.New("department still has users")
	}

	// Soft delete
	return repository.DeleteDepartment(id)
}
