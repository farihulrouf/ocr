package routes

import (
	"ocr-saas-backend/internal/handler"
	"ocr-saas-backend/internal/handler/categories"
	"ocr-saas-backend/internal/handler/dashboard"
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
	// PUBLIC AUTH ROUTES
	// =============================
	auth := v0.Group("/auth")
	auth.Post("/login", handler.Login)
	auth.Post("/refresh-token", handler.RefreshToken)
	// auth.Post("/forgot-password", handler.ForgotPassword)
	// auth.Post("/reset-password", handler.ResetPassword)
	// auth.Post("/verify-email", handler.VerifyEmail)

	// =============================
	// AUTH PROTECTED
	// =============================
	authProtected := auth.Group("/", middleware.Protected())
	authProtected.Get("/me", handler.GetProfile)
	authProtected.Put("/profile", handler.UpdateProfile)
	authProtected.Put("/password", handler.UpdatePassword)
	authProtected.Post("/logout", handler.Logout)

	// =============================
	// TENANT PROTECTED
	// =============================
	tenant := v0.Group("/tenant", middleware.Protected())

	tenant.Get("/info", tenants.GetTenantInfo)
	tenant.Put("/info", tenants.UpdateTenantInfo)
	tenant.Get("/settings", tenants.GetTenantSettings)

	tenant.Get("/subscription", tenants.GetTenantSubscription)
	tenant.Post("/subscription/upgrade", tenants.UpgradeSubscription)
	//SuperAdminOnly
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

	system.Get("/vendors", vendors.GetAllVendors)
	system.Post("/vendors", vendors.CreateVendor)
	system.Put("/vendors/:id", vendors.UpdateVendor)
	system.Delete("/vendors/:id", vendors.DeleteVendor)

	//finance := api.Group("/finance", middleware.Auth())

	//finance.Get("/categories", handler.ListCategories)
	//finance.Post("/categories", middleware.RoleAdmin(), handler.CreateCategory)
	//finance.Put("/categories/:id", middleware.RoleAdmin(), handler.UpdateCategory)
	//finance.Delete("/categories/:id", middleware.RoleAdmin(), handler.DeleteCategory)

	emprole := v0.Group("/emp", middleware.Protected(), middleware.EmployeeOnly())
	emprole.Get("/receipt", receipts.GetMyReceipts)

	emprole.Get("/receipt/:id", receipts.GetMyReceiptDetail)
	emprole.Post("/receipt/upload", ocr.UploadReceipt)
	emprole.Put("/receipt/:id", receipts.UpdateReceipt)
	//emprole.Put("/receipt/:id", handler.ConfirmReceipt)
	emprole.Delete("/receipt/:id", receipts.DeleteReceipt)
	emprole.Post("/receipt/:id/items", receipts.AddReceiptItem)
	emprole.Put("/receipt/items/:itemId", receipts.UpdateReceiptItem)

	//api.Post("/ocr/receipt", handler.UploadReceipt)

	// =============================
	// EMPLOYEE - EXPENSE REPORT
	// =============================
	//empReport := emprole.Group("/reports")
	emprole.Get("/reports/", reports.GetMyReports)
	emprole.Post("/reports/", reports.CreateReport)
	emprole.Put("/reports/:id", reports.UpdateReport)
	emprole.Post("/reports/:id/submit", reports.SubmitReport)
	// --- TAMBAHKAN INI UNTUK HAPUS STRUK DARI LAPORAN ---
	emprole.Delete("/reports/:id/receipts/:receiptId", reports.RemoveReceiptFromReport)
	emprole.Post("/reports/:id/receipts", reports.AddReceiptsToReport)
	emprole.Get("/reports/:id", reports.GetMyReportDetail)
	emprole.Get("/dashboard", dashboard.GetEmployeeDashboard)
	// OCR Upload

	//emprole.Post("/receipt/upload", ocr.UploadOCR)

	//app.Post("/v0/api/receipts/upload", receiptHandler.UploadOCR)

	//Get("/receipts", handler.GetMyReceipts)
	// =============================
	// USAGE STATS (ini yang kamu buat)
	// =============================
	//tenant.Get("/usage", handler.GetUsageStats) // GET /v0/api/tenant/usage
	manager := v0.Group("/manager", middleware.Protected(), middleware.TenantAdminOnly())
	manager.Get("/receipt", receipts.GetAllReceipts)
	manager.Get("/receipt/:id", receipts.GetReceiptDetail)
	manager.Put("/receipt/:id", receipts.ConfirmReceipt)
	//manager.Delete("/receipt/:id", handler.DeleteReceipt)
	manager.Post("/receipt/bulk/delete", receipts.BulkDeleteReceipts)
	manager.Post("/receipt/bulk/restore", receipts.BulkRestoreReceipts)
	manager.Post("/receipt/bulk/approve", receipts.BulkApproveReceipts)
	manager.Post("/receipt/bulk/reject", receipts.BulkRejectReceipts)
	manager.Post("/receipt/bulk/update-category", receipts.BulkUpdateReceiptCategory)
	manager.Post("/receipt/:id/items", receipts.AddReceiptItem)
	manager.Put("/receipt/items/:itemId", receipts.UpdateReceiptItem)
	manager.Delete("/receipt/items/:itemId", receipts.DeleteReceiptItem)
	manager.Get("/dashboard", dashboard.GetManagerDashboard)
	// =============================
	// MANAGER - REPORT APPROVAL
	// =============================
	manager.Get("/reports", reports.GetPendingReports)
	manager.Post("/reports/:id/approve", reports.ApproveReport)
	manager.Post("/reports/:id/reject", reports.RejectReport)

	// =============================
	// MANAGER - REPORT APPROVAL
	// =============================
	//managerReport := manager.Group("/reports")
	//manager.Get("/reports/", handler.GetPendingReports)
	//managerReport.Post("/:id/approve", handler.ApproveReport)
	//managerReport.Post("/:id/reject", handler.RejectReport)

}
