package categories

import (
	"ocr-saas-backend/internal/dto"
	"ocr-saas-backend/internal/service/categories"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func ListCategories(c *fiber.Ctx) error {
	tenantID := uuid.MustParse(c.Locals("tenant_id").(string))

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	resp, err := categories.ListCategories(tenantID, page, pageSize)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	return c.JSON(resp)
}

func CreateCategory(c *fiber.Ctx) error {
	tenantID := uuid.MustParse(c.Locals("tenant_id").(string))

	var req dto.AccountCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid body"})
	}

	if err := categories.CreateCategory(tenantID, req); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success"})
}

func UpdateCategory(c *fiber.Ctx) error {
	tenantID := uuid.MustParse(c.Locals("tenant_id").(string))
	id := uuid.MustParse(c.Params("id"))

	var req dto.AccountCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid body"})
	}

	if err := categories.UpdateCategory(tenantID, id, req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success"})
}

func DeleteCategory(c *fiber.Ctx) error {
	tenantID := uuid.MustParse(c.Locals("tenant_id").(string))
	id := uuid.MustParse(c.Params("id"))

	if err := categories.DeleteCategory(tenantID, id); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success"})
}
