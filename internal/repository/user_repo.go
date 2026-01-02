package repository

import (
	"ocr-saas-backend/configs"
	"ocr-saas-backend/internal/models"
)

func GetAllUsers(page, pageSize int, q, sort string) ([]map[string]interface{}, int64, error) {
	var users []models.User
	var total int64

	offset := (page - 1) * pageSize

	query := configs.DB.Model(&models.User{}).Preload("Department")

	// search
	if q != "" {
		query = query.Where("name ILIKE ? OR email ILIKE ?", "%"+q+"%", "%"+q+"%")
	}

	// count
	query.Count(&total)

	// sorting
	if sort != "" {
		query = query.Order(sort)
	} else {
		query = query.Order("created_at DESC")
	}

	// query data
	if err := query.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	// mapping output format
	result := make([]map[string]interface{}, 0)
	for _, u := range users {
		result = append(result, map[string]interface{}{
			"id":   u.ID,
			"name": u.Name,
			"role": u.Role,
			"dept": func() string {
				if u.Department != nil {
					return u.Department.Name
				}
				return ""
			}(),
		})
	}

	return result, total, nil
}
