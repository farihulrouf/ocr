package budgets

import (
	budgetService "ocr-saas-backend/internal/service/budgets"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// HandleSetBudget - Endpoint untuk Manager mengatur limit budget
func HandleSetBudget(c *fiber.Ctx) error {
	// Ambil metadata dari Middleware Auth
	tenantIDStr, _ := c.Locals("tenant_id").(string)
	userIDStr, _ := c.Locals("user_id").(string)

	type Request struct {
		Category string `json:"category"`
		Limit    int64  `json:"limit"`
		Month    int    `json:"month"`
		Year     int    `json:"year"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Input format is invalid",
		})
	}

	// Parsing UUID dengan aman
	tID, _ := uuid.Parse(tenantIDStr)
	uID, _ := uuid.Parse(userIDStr)

	err := budgetService.SetBudgetLimit(
		tID,
		uID,
		req.Category,
		req.Limit,
		req.Month,
		req.Year,
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Budget limit updated successfully",
	})
}

// ListBudgets - Mengambil history budget per tahun
func ListBudgets(c *fiber.Ctx) error {
	tenantIDStr, _ := c.Locals("tenant_id").(string)
	yearStr := c.Query("year")
	year, _ := strconv.Atoi(yearStr)

	tID, _ := uuid.Parse(tenantIDStr)
	data, err := budgetService.GetTenantBudgets(tID, year)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   data,
	})
}

// GetBudgetStats - INI ENDPOINT UTAMA DASHBOARD
func GetBudgetStats(c *fiber.Ctx) error {
	tenantIDStr, _ := c.Locals("tenant_id").(string)
	// Default kategori ke "General" jika tidak ada di query
	category := c.Query("category", "General")

	tID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid Tenant ID"})
	}

	// Panggil Service yang sudah kita pasangi Debug tadi
	stats, err := budgetService.CalculateBudgetStats(tID, category)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// Pastikan response data berisi 'stats' secara utuh
	// Di sinilah data 'estimated_end_month' dikirim ke Frontend
	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   stats,
	})
}

// GetFinanceSummary - Endpoint untuk Dashboard Global Finance
func GetFinanceSummary(c *fiber.Ctx) error {
	// 1. Ambil Tenant ID dari middleware auth
	tenantIDStr, _ := c.Locals("tenant_id").(string)
	tID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid Tenant ID",
		})
	}

	// 2. Ambil parameter Month & Year dari Query URL
	// Default ke bulan & tahun sekarang jika tidak dikirim dari Frontend
	now := time.Now()
	monthStr := c.Query("month", strconv.Itoa(int(now.Month())))
	yearStr := c.Query("year", strconv.Itoa(now.Year()))

	month, _ := strconv.Atoi(monthStr)
	year, _ := strconv.Atoi(yearStr)

	// 3. Panggil Service untuk hitung agregat (Summary)
	summary, err := budgetService.GetFinanceBudgetSummary(tID, month, year)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// 4. Set Header Charset agar karakter aman (untuk nama dept/manager jika ada)
	c.Set("Content-Type", "application/json; charset=utf-8")

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   summary,
	})
}
