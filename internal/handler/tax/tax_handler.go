package tax

import (
	"ocr-saas-backend/internal/service/tax"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetTaxRates(c *fiber.Ctx) error {
	tenantID := uuid.MustParse(c.Locals("tenant_id").(string))

	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	if pageSize < 1 {
		pageSize = 10
	}

	response, err := tax.GetTaxRates(
		tenantID, page, pageSize,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.JSON(response)
}

func CreateTaxRate(c *fiber.Ctx) error {
	tenantID := uuid.MustParse(c.Locals("tenant_id").(string))

	var body struct {
		Name       string `json:"name"`
		Percentage int    `json:"percentage"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid body"})
	}

	if err := tax.CreateTaxRate(
		tenantID,
		body.Name,
		body.Percentage,
	); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status": "success",
	})
}

func UpdateTaxRate(c *fiber.Ctx) error {
	tenantID := uuid.MustParse(c.Locals("tenant_id").(string))
	id := uuid.MustParse(c.Params("id"))

	var body struct {
		Name       string `json:"name"`
		Percentage int    `json:"percentage"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid body"})
	}

	if err := tax.UpdateTaxRate(
		tenantID,
		id,
		body.Name,
		body.Percentage,
	); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status": "success",
	})
}

func DeleteTaxRate(c *fiber.Ctx) error {
	tenantID := uuid.MustParse(c.Locals("tenant_id").(string))
	id := uuid.MustParse(c.Params("id"))

	if err := tax.DeleteTaxRate(
		tenantID,
		id,
	); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status": "success",
	})
}
