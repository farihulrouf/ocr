package budgets

import (
	budgetService "ocr-saas-backend/internal/service/budgets"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func HandleSetBudget(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(string)
	userID, _ := c.Locals("user_id").(string)

	type Request struct {
		Category string `json:"category"`
		Limit    int64  `json:"limit"`
		Month    int    `json:"month"`
		Year     int    `json:"year"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Format input salah"})
	}

	err := budgetService.SetBudgetLimit(
		uuid.MustParse(tenantID),
		uuid.MustParse(userID),
		req.Category,
		req.Limit,
		req.Month,
		req.Year,
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Budget berhasil diperbarui"})
}

func ListBudgets(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(string)
	yearStr := c.Query("year")
	year, _ := strconv.Atoi(yearStr)

	data, err := budgetService.GetTenantBudgets(uuid.MustParse(tenantID), year)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "data": data})
}

func GetBudgetStats(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(string)
	category := c.Query("category", "General")

	stats, err := budgetService.CalculateBudgetStats(uuid.MustParse(tenantID), category)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "data": stats})
}
