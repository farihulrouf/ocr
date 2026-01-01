package repository

import (
	"ocr-saas-backend/configs"
	"ocr-saas-backend/internal/models"
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
