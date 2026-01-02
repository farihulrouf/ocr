package handler

import (
	"ocr-saas-backend/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func ListUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	q := c.Query("q", "")
	sort := c.Query("sort", "")

	result, err := service.GetAllUsers(page, pageSize, q, sort)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}

	return c.Status(200).JSON(result)
}

func UserDetail(c *fiber.Ctx) error {
	tenantID, _ := uuid.Parse(c.Locals("tenant_id").(string))
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid user id"})
	}

	data, err := service.GetDetail(tenantID, userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "user not found"})
	}

	return c.JSON(data)
}

func UpdateUser(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid user id"})
	}

	tenantID := c.Locals("tenant_id").(uuid.UUID)

	var body struct {
		Role   string     `json:"role"`
		DeptID *uuid.UUID `json:"dept_id"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	resp, err := service.UpdateUser(tenantID, userID, body.Role, body.DeptID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(resp)
}
