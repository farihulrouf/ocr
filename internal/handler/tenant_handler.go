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

type UpdateTenantInfoReq struct {
	CompanyName string `json:"company_name"`
	TaxID       string `json:"tax_id"`
}

func UpdateTenantInfo(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id")
	if tenantID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req UpdateTenantInfoReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	updateData := map[string]interface{}{}
	if req.CompanyName != "" {
		updateData["name"] = req.CompanyName
	}
	if req.TaxID != "" {
		updateData["business_number"] = req.TaxID
	}

	if len(updateData) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No valid fields to update",
		})
	}

	if err := service.UpdateTenantInfo(tenantID.(string), updateData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update tenant info",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Info updated",
	})
}

func GetTenantSettings(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id")

	if tenantID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	settings, err := service.GetTenantSettings(tenantID.(string))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Settings not found",
		})
	}

	return c.JSON(fiber.Map{
		"currency":    settings.Currency,
		"date_format": settings.DateFormat,
		"ocr_auto":    settings.AutoOCR,
	})
}

func GetTenantSubscription(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id")
	if tenantID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	tenant, err := service.GetTenantSubscription(tenantID.(string))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Subscription not found",
		})
	}

	return c.JSON(fiber.Map{
		"plan_id":  tenant.SubscriptionPlanID,
		"status":   tenant.Status,
		"exp_date": nil, // karena model belum punya
	})
}

type UpgradeRequest struct {
	PlanID string `json:"plan_id"`
}

func UpgradeSubscription(c *fiber.Ctx) error {
	// ambil tenant dari middleware auth
	tenantID := c.Locals("tenant_id")
	if tenantID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var body UpgradeRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid body",
		})
	}

	checkoutURL, err := service.CreateUpgradeCheckoutURL(tenantID.(string), body.PlanID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Plan not found",
		})
	}

	return c.JSON(fiber.Map{
		"checkout_url": checkoutURL,
	})
}
