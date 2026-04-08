package reports

import (
	"fmt" // Pastikan import fmt ini ada
	"log"
	"ocr-saas-backend/internal/service/reports"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// HandleExport - Handler untuk POST /v0/api/manager/reports/export
func HandleExport(c *fiber.Ctx) error {
	log.Println("📩 [Handler] Menerima request export...")

	// 1. Ambil tenant_id & user_id dari context (interface{})
	tenantIDRaw := c.Locals("tenant_id")
	userIDRaw := c.Locals("user_id")

	if tenantIDRaw == nil || userIDRaw == nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Unauthorized: Tenant ID atau User ID tidak ditemukan di session",
		})
	}

	// 2. Konversi ke uuid.UUID secara aman (Menggunakan fmt.Sprintf yang asli)
	tenantID, err := uuid.Parse(fmt.Sprintf("%v", tenantIDRaw))
	if err != nil {
		log.Printf("❌ [Handler] Format TenantID salah: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Tenant ID format"})
	}

	userID, err := uuid.Parse(fmt.Sprintf("%v", userIDRaw))
	if err != nil {
		log.Printf("❌ [Handler] Format UserID salah: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid User ID format"})
	}

	// 3. Ambil status dari query
	status := c.Query("status", "APPROVED")

	log.Printf("🚀 [Handler] Memanggil service untuk Tenant: %s, User: %s, Status: %s",
		tenantID, userID, status)

	// 4. Panggil Service
	// Pastikan package internal/service/reports sudah di-import sebagai 'reports'
	err = reports.TriggerExpenseExport(c.Context(), tenantID, userID, status)
	if err != nil {
		log.Printf("❌ [Handler] Gagal memicu service: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal memproses export: " + err.Error(),
		})
	}

	// 5. Response Sukses
	return c.JSON(fiber.Map{
		"message": "Export sedang diproses, silakan cek riwayat download secara berkala",
		"details": fiber.Map{
			"status":    status,
			"tenant_id": tenantID,
		},
	})
}

// HandleGetExportLogs - GET /v0/api/manager/reports/export/logs
func HandleGetExportLogs(c *fiber.Ctx) error {
	tenantIDRaw := c.Locals("tenant_id")
	if tenantIDRaw == nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	tenantID, _ := uuid.Parse(fmt.Sprintf("%v", tenantIDRaw))

	// Ambil limit dari query param (misal: ?limit=1)
	limitStr := c.Query("limit", "10")
	limit, _ := strconv.Atoi(limitStr)

	// Panggil Service
	logs, err := reports.GetRecentExportLogs(c.Context(), tenantID, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil riwayat"})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   logs,
	})
}
