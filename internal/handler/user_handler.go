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

	tenantAny := c.Locals("tenant_id")
	if tenantAny == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	tenantIDStr, ok := tenantAny.(string)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "invalid tenant"})
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "invalid tenant uuid"})
	}

	result, err := service.GetAllUsers(
		tenantID,
		page,
		pageSize,
		q,
		sort,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
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
