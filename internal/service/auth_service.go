package service

import (
	"errors"
	"ocr-saas-backend/internal/models"
	"ocr-saas-backend/internal/repository"
	"ocr-saas-backend/pkg/utils"
)

type LoginResponse struct {
	Token        string      `json:"token"`
	RefreshToken string      `json:"refresh_token"`
	User         models.User `json:"user_obj"`
}

func Login(email, password string) (*LoginResponse, error) {
	// 1. Cari user berdasarkan email
	user, err := repository.FindUserByEmail(email)
	if err != nil {
		return nil, errors.New("email atau password salah")
	}

	// 2. Verifikasi Password
	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return nil, errors.New("email atau password salah")
	}

	// 3. Generate Tokens
	token, refresh, err := utils.GenerateToken(user.ID, user.TenantID, user.Role)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token:        token,
		RefreshToken: refresh,
		User:         *user,
	}, nil
}

// RESPONSE REFRESH TOKEN
type RefreshResponse struct {
	Token string `json:"token"`
}

func RefreshToken(refreshToken string) (*RefreshResponse, error) {
	claims, err := utils.ParseToken(refreshToken)

	if err != nil {
		return nil, errors.New("refresh token tidak valid")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("invalid refresh token claims")
	}

	user, err := repository.FindUserByID(userID)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	newToken, _, err := utils.GenerateToken(user.ID, user.TenantID, user.Role)
	if err != nil {
		return nil, errors.New("gagal membuat token baru")
	}

	return &RefreshResponse{Token: newToken}, nil
}

type ProfileResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
	Dept  string `json:"dept"`
}

func GetProfile(token string) (*ProfileResponse, error) {

	claims, err := utils.ParseToken(token)
	if err != nil {
		return nil, errors.New("token tidak valid")
	}

	userID := claims["user_id"].(string)

	user, err := repository.FindUserByID(userID)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	return &ProfileResponse{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
		Dept:  user.DepartmentID.String(),
	}, nil
}

func UpdateProfile(userID string, name string, avatar string) error {
	return repository.UpdateUserProfile(userID, name, avatar)
}

func UpdatePassword(userID string, oldPass string, newPass string) error {
	user, err := repository.FindUserByID(userID)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	if !utils.CheckPasswordHash(oldPass, user.PasswordHash) {
		return errors.New("password lama salah")
	}

	newHash, err := utils.HashPassword(newPass)
	if err != nil {
		return errors.New("gagal hash password")
	}

	return repository.UpdatePassword(userID, newHash)
}

func Logout(userID string) error {
	return repository.DeleteRefreshToken(userID)
}
