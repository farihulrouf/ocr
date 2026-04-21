package routes

import (
	"ocr-saas-backend/internal/handler"
	"ocr-saas-backend/internal/handler/budgets"
	"ocr-saas-backend/internal/handler/categories"
	"ocr-saas-backend/internal/handler/dashboard"
	"ocr-saas-backend/internal/handler/disbursement"
	"ocr-saas-backend/internal/handler/receipts"
	"ocr-saas-backend/internal/handler/tenants"

	"ocr-saas-backend/internal/handler/departments"
	"ocr-saas-backend/internal/handler/ocr"
	"ocr-saas-backend/internal/handler/payments"
	"ocr-saas-backend/internal/handler/reports"
	"ocr-saas-backend/internal/handler/tax"
	"ocr-saas-backend/internal/handler/vendors"
	"ocr-saas-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	v0 := app.Group("/v0/api")

	// =============================
	// AUTH
	// =============================
	auth := v0.Group("/auth")
	auth.Post("/login", handler.Login)
	auth.Post("/refresh-token", handler.RefreshToken)

	authProtected := auth.Group("/", middleware.Protected())
	authProtected.Get("/me", handler.GetProfile)
	authProtected.Put("/profile", handler.UpdateProfile)
	authProtected.Put("/password", handler.UpdatePassword)
	authProtected.Post("/logout", handler.Logout)

	// =============================
	// TENANT
	// =============================
	tenant := v0.Group("/tenant", middleware.Protected())
	tenant.Get("/info", tenants.GetTenantInfo)
	tenant.Put("/info", tenants.UpdateTenantInfo)
	tenant.Get("/settings", tenants.GetTenantSettings)
	tenant.Get("/subscription", tenants.GetTenantSubscription)
	tenant.Post("/subscription/upgrade", tenants.UpgradeSubscription)

	// =============================
	// SYSTEM (ADMIN ONLY)
	// =============================
	system := v0.Group("/system", middleware.Protected(), middleware.SuperAdminOnly())

	system.Get("/tenants", tenants.SystemListTenants)

	system.Get("/departments", departments.ListDepartments)
	system.Post("/departments", departments.CreateDepartment)
	system.Get("/departments/:id", departments.GetDepartmentDetailHandler)
	system.Put("/departments/:id", departments.UpdateDepartment)
	system.Delete("/departments/:id", departments.DeleteDepartment)

	system.Get("/org/users", handler.ListUsers)
	system.Get("/org/users/:id", handler.UserDetail)
	system.Put("/org/users/:id", handler.UpdateUser)

	system.Get("/categories", categories.ListCategories)
	system.Post("/categories", categories.CreateCategory)
	system.Put("/categories/:id", categories.UpdateCategory)
	system.Delete("/categories/:id", categories.DeleteCategory)

	system.Get("/tax", tax.GetTaxRates)
	system.Post("/tax", tax.CreateTaxRate)
	system.Put("/tax/:id", tax.UpdateTaxRate)
	system.Delete("/tax/:id", tax.DeleteTaxRate)

	system.Get("/payments", payments.GetAllPayments)
	system.Post("/payments", payments.CreatePayments)
	system.Put("/payments/:id", payments.UpdatePayments)
	system.Delete("/payments/:id", payments.DeletePayments)

	// =============================
	// EMPLOYEE
	// =============================
	emprole := v0.Group("/emp", middleware.Protected(), middleware.EmployeeOnly())

	emprole.Get("/receipt", receipts.GetMyReceipts)
	emprole.Get("/receipt/:id", receipts.GetMyReceiptDetail)
	emprole.Post("/receipt/upload", ocr.UploadReceipt)
	emprole.Put("/receipt/:id", receipts.UpdateReceipt)
	emprole.Delete("/receipt/:id", receipts.DeleteReceipt)

	emprole.Post("/receipt/:id/items", receipts.AddReceiptItem)
	emprole.Put("/receipt/items/:itemId", receipts.UpdateReceiptItem)
	emprole.Get("/receipt/:id/status", receipts.GetReceiptStatusHandler)

	// REPORTS (EMPLOYEE)
	emprole.Get("/reports", reports.GetMyReports)
	emprole.Post("/reports", reports.CreateReport)
	emprole.Put("/reports/:id", reports.UpdateReport)
	emprole.Post("/reports/:id/submit", reports.SubmitReport)
	emprole.Delete("/reports/:id/receipts/:receiptId", reports.RemoveReceiptFromReport)
	emprole.Post("/reports/:id/receipts", reports.AddReceiptsToReport)
	emprole.Get("/reports/:id", reports.GetMyReportDetail)

	emprole.Get("/dashboard", dashboard.GetEmployeeDashboard)

	// =============================
	// MANAGER
	// =============================
	manager := v0.Group("/manager", middleware.Protected(), middleware.TenantAdminOnly())

	manager.Get("/receipt", receipts.GetAllReceipts)
	manager.Get("/receipt/:id", receipts.GetReceiptDetail)
	manager.Put("/receipt/:id", receipts.ConfirmReceipt)

	manager.Post("/receipt/bulk/delete", receipts.BulkDeleteReceipts)
	manager.Post("/receipt/bulk/restore", receipts.BulkRestoreReceipts)
	manager.Post("/receipt/bulk/approve", receipts.BulkApproveReceipts)
	manager.Post("/receipt/bulk/reject", receipts.BulkRejectReceipts)
	manager.Post("/receipt/bulk/update-category", receipts.BulkUpdateReceiptCategory)

	manager.Post("/receipt/:id/items", receipts.AddReceiptItem)
	manager.Put("/receipt/items/:itemId", receipts.UpdateReceiptItem)
	manager.Delete("/receipt/items/:itemId", receipts.DeleteReceiptItem)

	manager.Get("/dashboard", dashboard.GetManagerDashboard)

	// REPORT APPROVAL
	manager.Get("/reports", reports.GetPendingReports)
	manager.Post("/reports/:id/approve", reports.ApproveReport)
	manager.Post("/reports/:id/reject", reports.RejectReport)

	manager.Post("/reports/export", reports.HandleExport)
	manager.Get("/reports/export/logs", reports.HandleGetExportLogs)

	// BUDGET
	manager.Post("/budget", budgets.HandleSetBudget)
	manager.Get("/budget/stats", budgets.GetBudgetStats)

	// =============================
	// SHARED REPORTS (MANAGER + FINANCE + ADMIN)
	// =============================
	reportsShared := v0.Group("/reports", middleware.Protected(), middleware.TenantAdminOnly())

	// ⚠️ STATIC HARUS DULU
	reportsShared.Get("/ready", reports.GetReadyToPayReports)

	// BARU DYNAMIC
	reportsShared.Get("/:id", reports.GetMyReportDetail)

	// =============================
	// FINANCE
	// =============================
	finance := v0.Group("/finance", middleware.Protected(), middleware.FinanceOnly())

	finance.Get("/dashboard", dashboard.GetAdminFinanceDashboard)

	// PAYMENT
	finance.Post("/pay", disbursement.ExecutePayment)

	// VENDORS
	finance.Get("/vendors", vendors.GetAllVendorsHandler)
	finance.Post("/vendors", vendors.CreateVendor)
	finance.Put("/vendors/:id", vendors.UpdateVendor)
	finance.Delete("/vendors/:id", vendors.DeleteVendor)

	// BUDGET
	finance.Get("/budget-summary", budgets.GetFinanceSummary)
	finance.Get("/budget", budgets.ListBudgets)
}
