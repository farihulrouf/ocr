package repository

import (
	"ocr-saas-backend/configs"
	"ocr-saas-backend/internal/models"

	"github.com/google/uuid"
)

func GetAllDepartments(page, pageSize int, q, sort string) ([]models.Department, int64, error) {
	var departments []models.Department
	var total int64

	db := configs.DB.Model(&models.Department{})

	// Search
	if q != "" {
		like := "%" + q + "%"
		db = db.Where("name LIKE ?", like)
	}

	// Count total rows
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Sorting
	if sort != "" {
		db = db.Order(sort)
	} else {
		db = db.Order("created_at DESC")
	}

	offset := (page - 1) * pageSize

	// Query data + pagination
	err := db.
		Limit(pageSize).
		Offset(offset).
		Find(&departments).Error

	return departments, total, err
}

func CreateDepartment(dept *models.Department) error {
	return configs.DB.Create(dept).Error
}

func GetDepartmentByID(id uuid.UUID) (*models.Department, error) {
	var dept models.Department

	err := configs.DB.
		First(&dept, "id = ?", id).Error

	if err != nil {
		return nil, err
	}

	return &dept, nil
}

func GetUsersByDepartment(deptID uuid.UUID) ([]models.User, error) {
	var users []models.User

	err := configs.DB.
		Where("department_id = ?", deptID).
		Find(&users).Error

	return users, err
}

func UpdateDepartmentByID(id uuid.UUID, name string) error {
	return configs.DB.Model(&models.Department{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"name": name,
		}).Error
}

func DeleteDepartment(id uuid.UUID) error {
	return configs.DB.Where("id = ?", id).Delete(&models.Department{}).Error
}
