package handler

import (
	"ocr-saas-backend/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetMyReceipts(c *fiber.Ctx) error {
	tenantID, ok := c.Locals("tenant_id").(string)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "unauthorized"})
	}

	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "unauthorized"})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	if pageSize < 1 {
		pageSize = 10
	}

	q := c.Query("q")
	status := c.Query("status")
	sort := c.Query("sort")

	response, err := service.GetMyReceipts(
		uuid.MustParse(tenantID),
		uuid.MustParse(userID),
		page,
		pageSize,
		q,
		status,
		sort,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.JSON(response)
}

func GetAllReceipts(c *fiber.Ctx) error {
	tenantID, ok := c.Locals("tenant_id").(string)
	if !ok {
		return c.Status(401).JSON(fiber.Map{
			"status":  "error",
			"message": "unauthorized",
		})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	if pageSize < 1 {
		pageSize = 10
	}

	q := c.Query("q")
	status := c.Query("status")
	sort := c.Query("sort")

	response, err := service.GetAllReceipts(
		uuid.MustParse(tenantID),
		page,
		pageSize,
		q,
		status,
		sort,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.JSON(response)
}
