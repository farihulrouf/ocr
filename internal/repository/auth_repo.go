package repository

import (
	"ocr-saas-backend/configs"
	"ocr-saas-backend/internal/models"
)

func FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	// Kita ambil data user beserta info Tenant-nya
	err := configs.DB.Preload("Tenant").Where("email = ?", email).First(&user).Error
	return &user, err
}

func FindUserByID(id string) (*models.User, error) {
	var user models.User
	err := configs.DB.Preload("Tenant").Where("id = ?", id).First(&user).Error
	return &user, err
}

func UpdateUserProfile(userID string, name string, avatar string) error {
	return configs.DB.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"name":   name,
			"avatar": avatar,
		}).Error
}

func UpdatePassword(userID string, newHash string) error {
	return configs.DB.Model(&models.User{}).
		Where("id = ?", userID).
		Update("password_hash", newHash).Error
}
