package reports

import (
	svc "ocr-saas-backend/internal/service/reports"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetMyReports(c *fiber.Ctx) error {
	// 1. Ambil data dari Middleware Protected
	tenantID := uuid.MustParse(c.Locals("tenant_id").(string))
	userID := uuid.MustParse(c.Locals("user_id").(string))
	role := c.Locals("role").(string) // <--- Ambil Role dari JWT

	// 2. Ambil Query Params
	status := c.Query("status", "")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	// 3. Panggil Service dengan menyertakan ROLE (6 Argumen)
	data, total, err := svc.GetMyReports(
		tenantID, userID, page, pageSize, status, role,
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
	newReport, err := svc.CreateReport(tenantID, userID, body.Title)
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

	if err := svc.SubmitReport(
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

	if err := svc.UpdateReport(
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

	data, total, err := svc.GetPendingReports(
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
	// 1. Ambil Tenant ID dari context (Middleware)
	tenantID := uuid.MustParse(c.Locals("tenant_id").(string))

	// 2. AMBIL USER ID (Manajer) dari context (Middleware)
	// Pastikan key-nya sesuai dengan yang ada di middleware Auth Mas (misal: "user_id")
	managerID := uuid.MustParse(c.Locals("user_id").(string))

	// 3. Ambil Report ID dari URL Parameter
	reportID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "ID laporan tidak valid"})
	}

	// 4. PANGGIL SERVICE dengan 3 PARAMETER (tenantID, reportID, managerID)
	if err := svc.ApproveReport(tenantID, reportID, managerID); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status":      "success",
		"message":     "Laporan berhasil disetujui",
		"approved_by": managerID, // Opsional: kirim balik ID approver-nya
	})
}

func RejectReport(c *fiber.Ctx) error {
	// 1. Ambil TenantID dan ManagerID dari JWT Locals
	tenantID := uuid.MustParse(c.Locals("tenant_id").(string))
	managerID := uuid.MustParse(c.Locals("user_id").(string)) // <-- Ambil ini

	// 2. Ambil ReportID dari URL Params
	reportID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "ID laporan tidak valid"})
	}

	// 3. Panggil service dengan 3 parameter
	if err := svc.RejectReport(tenantID, reportID, managerID); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "rejected_by": managerID})
}

func GetMyReportDetail(c *fiber.Ctx) error {
	tenantID := uuid.MustParse(c.Locals("tenant_id").(string))
	userID := uuid.MustParse(c.Locals("user_id").(string))
	role := c.Locals("role").(string) // Ambil role dari middleware
	reportID := uuid.MustParse(c.Params("id"))

	data, err := svc.GetMyReportDetail(
		tenantID, userID, reportID, role,
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

	err := svc.AddReceiptsToReport(tenantID, reportID, ids)
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

	err := svc.RemoveReceiptFromReport(tenantID, reportID, receiptID)
	if err != nil {
		return fiber.NewError(500, err.Error())
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Receipt removed from report"})
}

// Handler: GET /v0/api/finance/reports/ready
func GetReadyToPayReports(c *fiber.Ctx) error {

	// =============================
	// PARSE TENANT ID (SAFE)
	// =============================
	tenantIDStr, ok := c.Locals("tenant_id").(string)
	if !ok {
		return c.Status(401).JSON(fiber.Map{
			"message": "Invalid tenant",
		})
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid tenant_id format",
		})
	}

	// =============================
	// ROLE
	// =============================
	role, ok := c.Locals("role").(string)
	if !ok {
		return c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	// =============================
	// ROLE ACCESS CONTROL (STRICT)
	// =============================
	// ROLE ACCESS CONTROL (FIX)
	if role != "FINANCE" && role != "ADMIN" && role != "MANAGER" {
		return c.Status(403).JSON(fiber.Map{
			"message": "Akses ditolak",
		})
	}

	// =============================
	// PAGINATION (SAFE PARSE)
	// =============================
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.Query("page_size", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// =============================
	// STATUS FILTER
	// =============================
	status := c.Query("status", "APPROVED")

	// 🔥 RULE KHUSUS FINANCE
	if role == "FINANCE" {
		if status != "APPROVED" && status != "PAID" && status != "REJECTED" {
			status = "APPROVED"
		}
	}

	// 🔥 OPTIONAL: MANAGER bebas lihat semua
	// kalau mau batasi juga, tinggal tambahin rule di sini

	// =============================
	// SERVICE CALL
	// =============================
	data, total, err := svc.GetReadyToPayReports(tenantID, page, pageSize, status)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// =============================
	// RESPONSE
	// =============================
	return c.JSON(fiber.Map{
		"status": "success",
		"data":   data,
		"meta": fiber.Map{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
			"status":    status,
			"user_role": role,
		},
	})
}

func BulkApproveReports(c *fiber.Ctx) error {
	// =============================
	// AUTH CONTEXT
	// =============================
	tenantID := uuid.MustParse(c.Locals("tenant_id").(string))
	managerID := uuid.MustParse(c.Locals("user_id").(string))
	role := c.Locals("role").(string)

	// =============================
	// ROLE GUARD
	// =============================
	if role != "MANAGER" && role != "ADMIN" {
		return c.Status(403).JSON(fiber.Map{
			"message": "Akses ditolak",
		})
	}

	// =============================
	// REQUEST BODY
	// =============================
	var body struct {
		ReportIDs []string `json:"report_ids"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "invalid body",
		})
	}

	if len(body.ReportIDs) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"message": "report_ids is required",
		})
	}

	// 🔥 LIMIT (best practice)
	if len(body.ReportIDs) > 50 {
		return c.Status(400).JSON(fiber.Map{
			"message": "max 50 reports per request",
		})
	}

	// =============================
	// PARSE UUIDS
	// =============================
	var reportIDs []uuid.UUID
	for _, idStr := range body.ReportIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"message": "invalid report id: " + idStr,
			})
		}
		reportIDs = append(reportIDs, id)
	}

	// =============================
	// SERVICE CALL
	// =============================
	if err := svc.BulkApproveReports(tenantID, reportIDs, managerID); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// =============================
	// RESPONSE
	// =============================
	return c.JSON(fiber.Map{
		"status":        "success",
		"approved_by":   managerID,
		"total_reports": len(reportIDs),
	})
}

func BulkRejectReports(c *fiber.Ctx) error {
	// =============================
	// AUTH
	// =============================
	tenantID := uuid.MustParse(c.Locals("tenant_id").(string))
	managerID := uuid.MustParse(c.Locals("user_id").(string))
	role := c.Locals("role").(string)

	// =============================
	// ROLE GUARD
	// =============================
	if role != "MANAGER" && role != "ADMIN" {
		return c.Status(403).JSON(fiber.Map{
			"message": "Akses ditolak",
		})
	}

	// =============================
	// BODY
	// =============================
	var body struct {
		ReportIDs []string `json:"report_ids"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "invalid body",
		})
	}

	if len(body.ReportIDs) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"message": "report_ids is required",
		})
	}

	if len(body.ReportIDs) > 50 {
		return c.Status(400).JSON(fiber.Map{
			"message": "max 50 reports per request",
		})
	}

	// =============================
	// PARSE UUID
	// =============================
	var reportIDs []uuid.UUID
	for _, idStr := range body.ReportIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"message": "invalid report id: " + idStr,
			})
		}
		reportIDs = append(reportIDs, id)
	}

	// =============================
	// SERVICE
	// =============================
	if err := svc.BulkRejectReports(tenantID, reportIDs, managerID); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":        "success",
		"rejected_by":   managerID,
		"total_reports": len(reportIDs),
	})
}

func BulkPayReports(c *fiber.Ctx) error {
	// =============================
	// AUTH
	// =============================
	tenantID := uuid.MustParse(c.Locals("tenant_id").(string))
	financeID := uuid.MustParse(c.Locals("user_id").(string))
	role := c.Locals("role").(string)

	// =============================
	// ROLE GUARD
	// =============================
	if role != "FINANCE" && role != "ADMIN" {
		return c.Status(403).JSON(fiber.Map{
			"message": "Akses ditolak",
		})
	}

	// =============================
	// BODY
	// =============================
	var body struct {
		ReportIDs []string `json:"report_ids"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "invalid body",
		})
	}

	if len(body.ReportIDs) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"message": "report_ids is required",
		})
	}

	if len(body.ReportIDs) > 50 {
		return c.Status(400).JSON(fiber.Map{
			"message": "max 50 reports per request",
		})
	}

	// =============================
	// PARSE UUID
	// =============================
	var reportIDs []uuid.UUID
	for _, idStr := range body.ReportIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"message": "invalid report id: " + idStr,
			})
		}
		reportIDs = append(reportIDs, id)
	}

	// =============================
	// SERVICE
	// =============================
	if err := svc.BulkPayReports(tenantID, reportIDs, financeID); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// =============================
	// RESPONSE
	// =============================
	return c.JSON(fiber.Map{
		"status":        "success",
		"paid_by":       financeID,
		"total_reports": len(reportIDs),
	})
}

func BulkFailPaymentReports(c *fiber.Ctx) error {
	// =============================
	// AUTH
	// =============================
	tenantID := uuid.MustParse(c.Locals("tenant_id").(string))
	financeID := uuid.MustParse(c.Locals("user_id").(string))
	role := c.Locals("role").(string)

	// =============================
	// ROLE GUARD
	// =============================
	if role != "FINANCE" && role != "ADMIN" {
		return c.Status(403).JSON(fiber.Map{
			"message": "Akses ditolak",
		})
	}

	// =============================
	// BODY
	// =============================
	var body struct {
		ReportIDs []string `json:"report_ids"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "invalid body",
		})
	}

	if len(body.ReportIDs) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"message": "report_ids is required",
		})
	}

	if len(body.ReportIDs) > 50 {
		return c.Status(400).JSON(fiber.Map{
			"message": "max 50 reports per request",
		})
	}

	// =============================
	// PARSE UUID
	// =============================
	var reportIDs []uuid.UUID
	for _, idStr := range body.ReportIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"message": "invalid report id: " + idStr,
			})
		}
		reportIDs = append(reportIDs, id)
	}

	// =============================
	// SERVICE
	// =============================
	if err := svc.BulkFailPaymentReports(tenantID, reportIDs, financeID); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// =============================
	// RESPONSE
	// =============================
	return c.JSON(fiber.Map{
		"status":        "success",
		"failed_by":     financeID,
		"total_reports": len(reportIDs),
	})
}
