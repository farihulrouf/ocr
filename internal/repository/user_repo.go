package repository

import (
	"ocr-saas-backend/configs"
	"ocr-saas-backend/internal/models"

	"github.com/google/uuid"
)

func GetAllUsers(
	tenantID uuid.UUID, page, pageSize int, q, sort string) ([]map[string]interface{}, int64, error) {
	var users []models.User
	var total int64

	offset := (page - 1) * pageSize

	query := configs.DB.
		Model(&models.User{}).
		Preload("Department").
		Where("tenant_id = ?", tenantID) // 🔒 PENTING

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

func GetUserDetail(tenantID, userID uuid.UUID) (map[string]interface{}, error) {
	var user models.User

	err := configs.DB.
		Where("id = ? AND tenant_id = ?", userID, tenantID).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	// ambil history (AuditTrail)
	var history []models.AuditTrail
	configs.DB.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&history)

	// mapping output
	result := map[string]interface{}{
		"id":      user.ID,
		"name":    user.Name,
		"history": history,
	}

	return result, nil
}

func UpdateUser(tenantID, userID uuid.UUID, role string, deptID *uuid.UUID) error {
	updateData := map[string]interface{}{
		"role": role,
	}

	// department_id boleh null
	if deptID != nil {
		updateData["department_id"] = *deptID
	} else {
		updateData["department_id"] = nil
	}

	return configs.DB.
		Model(&models.User{}).
		Where("id = ? AND tenant_id = ?", userID, tenantID).
		Updates(updateData).Error
}
