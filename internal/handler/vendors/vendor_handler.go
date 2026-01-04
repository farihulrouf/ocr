package vendors

import (
	"ocr-saas-backend/internal/service/vendors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetAllVendors(c *fiber.Ctx) error {
	tenantID := uuid.MustParse(c.Locals("tenant_id").(string))

	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	if pageSize < 1 {
		pageSize = 10
	}

	response, err := vendors.GetAllVendor(
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

func CreateVendor(c *fiber.Ctx) error {
	tenantID := uuid.MustParse(c.Locals("tenant_id").(string))

	var body struct {
		Name      string `json:"name"`
		TaxNumber string `json:"tax_number"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid body"})
	}

	if err := vendors.CreateVendor(
		tenantID,
		body.Name,
		body.TaxNumber,
	); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status": "success",
	})
}

func UpdateVendor(c *fiber.Ctx) error {
	tenantID := uuid.MustParse(c.Locals("tenant_id").(string))
	id := uuid.MustParse(c.Params("id"))

	var body struct {
		Name      string `json:"name"`
		TaxNumber string `json:"tax_number"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid body"})
	}

	if err := vendors.UpdateVendor(
		tenantID,
		id,
		body.Name,
		body.TaxNumber,
	); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status": "success",
	})
}

func DeleteVendor(c *fiber.Ctx) error {
	tenantID := uuid.MustParse(c.Locals("tenant_id").(string))
	id := uuid.MustParse(c.Params("id"))

	if err := vendors.DeleteVendor(
		tenantID,
		id,
	); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status": "success",
	})
}
