package utils

import (
	"errors" // ← ini yang kurang
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword membuat hash dari password plain text
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12) // cost 12 aman
	return string(bytes), err
}

// CheckPasswordHash membandingkan password input dengan hash di DB
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateToken membuat Access Token & Refresh Token
func GenerateToken(userID uuid.UUID, tenantID uuid.UUID, role string) (string, string, error) {
	secret := os.Getenv("JWT_SECRET")

	// 1. Access Token (Berlaku 15 Menit/1 Jam)
	accessClaims := jwt.MapClaims{
		"user_id":   userID.String(),
		"tenant_id": tenantID.String(),
		"role":      role,
		"exp":       time.Now().Add(time.Hour * 24).Unix(), // Kita set 24 jam dulu untuk dev
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(secret))

	// 2. Refresh Token (Berlaku 7 Hari)
	refreshClaims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	}
	refreshToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(secret))

	return accessToken, refreshToken, err
}

func ParseToken(tokenString string) (jwt.MapClaims, error) {
	secret := os.Getenv("JWT_SECRET")

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// Pastikan HMAC SHA256
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("metode signing token tidak valid")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, errors.New("token tidak valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token tidak valid")
	}

	return claims, nil
}
