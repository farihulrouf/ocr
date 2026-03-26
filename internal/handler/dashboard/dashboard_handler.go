package dashboard

import (
	"ocr-saas-backend/internal/service/dashboard"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetEmployeeDashboard(c *fiber.Ctx) error {
	// Ambil dari Middleware Protected
	tenantID := uuid.MustParse(c.Locals("tenant_id").(string))
	userID := uuid.MustParse(c.Locals("user_id").(string))

	data, err := dashboard.GetEmployeeDashboardData(tenantID, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(data)
}
