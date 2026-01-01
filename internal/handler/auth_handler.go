package handler

import (
	"fmt"
	"ocr-saas-backend/internal/service"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx) error {
	type LoginRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Format request salah"})
	}

	// Panggil service
	result, err := service.Login(req.Email, req.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   result,
	})
}

func RefreshToken(c *fiber.Ctx) error {
	type RefreshTokenRequest struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	var req RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Format request salah",
		})
	}
	fmt.Println("cek", req.RefreshToken)
	result, err := service.RefreshToken(req.RefreshToken)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   result,
	})
}

func GetProfile(c *fiber.Ctx) error {

	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{
			"status":  "error",
			"message": "Authorization header missing",
		})
	}

	tokenStr := strings.Replace(authHeader, "Bearer ", "", 1)

	result, err := service.GetProfile(tokenStr)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   result,
	})
}
