package handler

import (
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

func UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id") // Dari middleware JWT kamu (belum kita buat)

	type Req struct {
		Name   string `json:"name"`
		Avatar string `json:"avatar"`
	}

	var req Req
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Format request salah",
		})
	}

	if err := service.UpdateProfile(userID.(string), req.Name, req.Avatar); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Profile updated",
	})
}

type UpdatePasswordRequest struct {
	OldPass string `json:"old_pass"`
	NewPass string `json:"new_pass"`
}

func UpdatePassword(c *fiber.Ctx) error {
	var body UpdatePasswordRequest

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	userID := c.Locals("user_id").(string)

	err := service.UpdatePassword(userID, body.OldPass, body.NewPass)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Password updated",
		"status":  "success",
	})
}

func Logout(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	err := service.Logout(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Session cleared",
	})
}
