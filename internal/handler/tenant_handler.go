package handler

import (
	"ocr-saas-backend/internal/service"

	"github.com/gofiber/fiber/v2"
)

func GetTenantInfo(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id")

	if tenantID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	tenant, err := service.GetTenantInfo(tenantID.(string))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Tenant not found",
		})
	}

	return c.JSON(fiber.Map{
		"company_name": tenant.Name,
		"tax_id":       tenant.BusinessNumber,
		"subdomain":    tenant.Subdomain,
		"status":       tenant.Status,
		"number":       tenant.BusinessNumber,
		"address":      "", // belum ada field
	})
}
