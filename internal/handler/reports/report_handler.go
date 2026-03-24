package reports

import (
	"ocr-saas-backend/internal/service/reports"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetMyReports(c *fiber.Ctx) error {
	tenantID := uuid.MustParse(c.Locals("tenant_id").(string))
	userID := uuid.MustParse(c.Locals("user_id").(string))

	// ✅ Ambil status dari query param: /api/emp/reports?status=DRAFT
	status := c.Query("status", "")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	data, total, err := reports.GetMyReports(
		tenantID, userID, page, pageSize, status,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   data,
		"meta": fiber.Map{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
			"status":    status, // Kembalikan status di meta untuk tracing
		},
	})
}

func CreateReport(c *fiber.Ctx) error {
	tenantID := uuid.MustParse(c.Locals("tenant_id").(string))
	userID := uuid.MustParse(c.Locals("user_id").(string))

	var body struct {
		Title string `json:"title"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid body"})
	}

	// Eksekusi create
	newReport, err := reports.CreateReport(tenantID, userID, body.Title)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	// RESPONSE KE FRONTEND
	// Frontend akan menerima: { "status": "success", "id": "uuid-laporan-baru" }
	return c.JSON(fiber.Map{
		"status": "success",
		"id":     newReport.ID,
	})
}

func SubmitReport(c *fiber.Ctx) error {
	tenantID := uuid.MustParse(c.Locals("tenant_id").(string))
	userID := uuid.MustParse(c.Locals("user_id").(string))
	reportID := uuid.MustParse(c.Params("id"))

	if err := reports.SubmitReport(
		tenantID, userID, reportID,
	); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success"})
}

func UpdateReport(c *fiber.Ctx) error {
	tenantID := uuid.MustParse(c.Locals("tenant_id").(string))
	userID := uuid.MustParse(c.Locals("user_id").(string))
	reportID := uuid.MustParse(c.Params("id"))

	var body struct {
		Title string `json:"title"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid body"})
	}

	if err := reports.UpdateReport(
		tenantID, userID, reportID, body.Title,
	); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success"})
}

func GetPendingReports(c *fiber.Ctx) error {
	tenantID := uuid.MustParse(c.Locals("tenant_id").(string))

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	data, total, err := reports.GetPendingReports(
		tenantID, page, pageSize,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   data,
		"meta": fiber.Map{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

func ApproveReport(c *fiber.Ctx) error {
	tenantID := uuid.MustParse(c.Locals("tenant_id").(string))
	reportID := uuid.MustParse(c.Params("id"))

	if err := reports.ApproveReport(tenantID, reportID); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success"})
}

func RejectReport(c *fiber.Ctx) error {
	tenantID := uuid.MustParse(c.Locals("tenant_id").(string))
	reportID := uuid.MustParse(c.Params("id"))

	if err := reports.RejectReport(tenantID, reportID); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success"})
}

func GetMyReportDetail(c *fiber.Ctx) error {
	tenantID := uuid.MustParse(c.Locals("tenant_id").(string))
	userID := uuid.MustParse(c.Locals("user_id").(string))
	reportID := uuid.MustParse(c.Params("id"))

	data, err := reports.GetMyReportDetail(
		tenantID, userID, reportID,
	)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   data,
	})
}

// POST /emp/reports/:id/receipts
func AddReceiptsToReport(c *fiber.Ctx) error {
	reportID := uuid.MustParse(c.Params("id"))
	tenantID := uuid.MustParse(c.Locals("tenant_id").(string))

	var req struct {
		ReceiptIDs []string `json:"receipt_ids"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(400, "invalid body")
	}

	ids := make([]uuid.UUID, len(req.ReceiptIDs))
	for i, s := range req.ReceiptIDs {
		id, err := uuid.Parse(s)
		if err != nil {
			return fiber.NewError(400, "invalid receipt id")
		}
		ids[i] = id
	}

	err := reports.AddReceiptsToReport(tenantID, reportID, ids)
	if err != nil {
		return fiber.NewError(500, err.Error())
	}

	return c.JSON(fiber.Map{"status": "success"})
}

// DELETE /emp/reports/:id/receipts/:receiptId
func RemoveReceiptFromReport(c *fiber.Ctx) error {
	reportID := uuid.MustParse(c.Params("id"))
	receiptID := uuid.MustParse(c.Params("receiptId"))
	tenantID := uuid.MustParse(c.Locals("tenant_id").(string))

	err := reports.RemoveReceiptFromReport(tenantID, reportID, receiptID)
	if err != nil {
		return fiber.NewError(500, err.Error())
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Receipt removed from report"})
}
