package handler

import (
	"ocr-saas-backend/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func ListDepartments(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	q := c.Query("q", "")
	sort := c.Query("sort", "")

	result, err := service.GetAllDepartments(page, pageSize, q, sort)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}

	return c.Status(200).JSON(result)
}

type CreateDeptRequest struct {
	Name string `json:"name"`
}

func CreateDepartment(c *fiber.Ctx) error {
	var body CreateDeptRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if body.Name == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Name is required",
		})
	}

	tenantID, _ := uuid.Parse(c.Locals("tenant_id").(string))

	dept, err := service.CreateDepartment(tenantID, body.Name)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create department",
		})
	}

	return c.JSON(fiber.Map{
		"id":   dept.ID,
		"name": dept.Name,
	})
}
