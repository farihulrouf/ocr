package vendors

import (
	"ocr-saas-backend/internal/service/vendors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetAllVendorsHandler(c *fiber.Ctx) error {
	// 1. Ambil Tenant ID dari Middleware (Locals)
	// Pastikan middleware auth kamu sudah set c.Locals("tenant_id", "uuid-string")
	rawTenantID := c.Locals("tenant_id")
	if rawTenantID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized: Tenant ID missing in context",
		})
	}

	tenantID, err := uuid.Parse(rawTenantID.(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid Tenant ID format",
		})
	}

	// 2. Ambil Query Parameters
	search := c.Query("search", "")
	category := c.Query("category", "")

	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	if pageSize < 1 {
		pageSize = 10
	}

	// 3. Panggil Service
	response, err := vendors.GetAllVendorService(
		tenantID, search, category, page, pageSize,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch vendors: " + err.Error(),
		})
	}

	// 4. Return JSON
	return c.JSON(response)
}

func CreateVendor(c *fiber.Ctx) error {
	// Ambil tenant_id dari middleware auth
	tenantID := uuid.MustParse(c.Locals("tenant_id").(string))

	var body struct {
		Name        string `json:"name"`
		DisplayName string `json:"display_name"`
		TaxNumber   string `json:"tax_number"`
		Category    string `json:"category"`
		Aliases     string `json:"aliases"`
		IsGlobal    bool   `json:"is_global"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid request body"})
	}

	// Panggil service create dengan parameter lengkap
	err := vendors.CreateVendor(
		tenantID,
		body.Name,
		body.DisplayName,
		body.TaxNumber,
		body.Category,
		body.Aliases,
		body.IsGlobal,
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "vendor created successfully",
	})
}

func UpdateVendor(c *fiber.Ctx) error {
	// Ambil tenant_id dari middleware auth
	tenantID := uuid.MustParse(c.Locals("tenant_id").(string))
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid vendor id"})
	}

	// Definisikan body sesuai dengan yang dikirim FE
	var body struct {
		Name        string `json:"name"`
		DisplayName string `json:"display_name"`
		TaxNumber   string `json:"tax_number"`
		Category    string `json:"category"`
		Aliases     string `json:"aliases"`
		IsGlobal    bool   `json:"is_global"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid request body"})
	}

	// Panggil service update
	err = vendors.UpdateVendorService(
		tenantID,
		id,
		body.Name,
		body.DisplayName,
		body.TaxNumber,
		body.Category,
		body.Aliases,
		body.IsGlobal,
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "vendor updated successfully",
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
