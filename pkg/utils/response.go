package utils

import "github.com/gofiber/fiber/v2"

// SuccessResponse untuk format sukses
func SuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": message,
		"data":    data,
	})
}

// ErrorResponse untuk format error
func ErrorResponse(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"status":  "error",
		"message": message,
	})
}