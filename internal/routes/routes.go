package routes

import (
	"ocr-saas-backend/internal/handler"
	"ocr-saas-backend/internal/handler/categories"
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

	tenant.Get("/info", handler.GetTenantInfo)
	tenant.Put("/info", handler.UpdateTenantInfo)
	tenant.Get("/settings", handler.GetTenantSettings)

	tenant.Get("/subscription", handler.GetTenantSubscription)
	tenant.Post("/subscription/upgrade", handler.UpgradeSubscription)
	//SuperAdminOnly
	system := v0.Group("/system", middleware.Protected(), middleware.SuperAdminOnly())
	system.Get("/tenants", handler.SystemListTenants)
	system.Get("/departments", handler.ListDepartments)
	system.Post("/departments", handler.CreateDepartment)
	system.Get("/departments/:id", handler.GetDepartmentDetailHandler)
	system.Put("/departments/:id", handler.UpdateDepartment)
	system.Delete("/departments/:id", handler.DeleteDepartment)

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
	emprole.Get("/receipt", handler.GetMyReceipts)

	emprole.Get("/receipt/:id", handler.GetMyReceiptDetail)
	emprole.Post("/receipt/upload", ocr.UploadReceipt)
	emprole.Put("/receipt/:id", handler.UpdateReceipt)
	//emprole.Put("/receipt/:id", handler.ConfirmReceipt)
	emprole.Delete("/receipt/:id", handler.DeleteReceipt)
	emprole.Post("/receipt/:id/items", handler.AddReceiptItem)
	emprole.Put("/receipt/items/:itemId", handler.UpdateReceiptItem)

	//api.Post("/ocr/receipt", handler.UploadReceipt)

	// =============================
	// EMPLOYEE - EXPENSE REPORT
	// =============================
	//empReport := emprole.Group("/reports")
	emprole.Get("/reports/", reports.GetMyReports)
	emprole.Post("/reports/", reports.CreateReport)
	emprole.Put("/reports/:id", reports.UpdateReport)
	emprole.Post("/reports/:id/submit", reports.SubmitReport)

	emprole.Post("/reports/:id/receipts", reports.AddReceiptsToReport)
	emprole.Get("/reports/:id", reports.GetMyReportDetail)

	// OCR Upload

	//emprole.Post("/receipt/upload", ocr.UploadOCR)

	//app.Post("/v0/api/receipts/upload", receiptHandler.UploadOCR)

	//Get("/receipts", handler.GetMyReceipts)
	// =============================
	// USAGE STATS (ini yang kamu buat)
	// =============================
	//tenant.Get("/usage", handler.GetUsageStats) // GET /v0/api/tenant/usage
	manager := v0.Group("/manager", middleware.Protected(), middleware.TenantAdminOnly())
	manager.Get("/receipt", handler.GetAllReceipts)
	manager.Get("/receipt/:id", handler.GetReceiptDetail)
	manager.Put("/receipt/:id", handler.ConfirmReceipt)
	//manager.Delete("/receipt/:id", handler.DeleteReceipt)
	manager.Post("/receipt/bulk/delete", handler.BulkDeleteReceipts)
	manager.Post("/receipt/bulk/restore", handler.BulkRestoreReceipts)
	manager.Post("/receipt/bulk/approve", handler.BulkApproveReceipts)
	manager.Post("/receipt/bulk/reject", handler.BulkRejectReceipts)
	manager.Post("/receipt/bulk/update-category", handler.BulkUpdateReceiptCategory)
	manager.Post("/receipt/:id/items", handler.AddReceiptItem)
	manager.Put("/receipt/items/:itemId", handler.UpdateReceiptItem)
	manager.Delete("/receipt/items/:itemId", handler.DeleteReceiptItem)
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

/*
package routes

import (
	"ocr-saas-backend/internal/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Base API Group
	v0 := app.Group("/v0/api")

	// =========================================================================
	// 1. AUTHENTICATION & SECURITY (10 API)
	// =========================================================================
	auth := v0.Group("/auth")
	auth.Post("/login", func(c *fiber.Ctx) error { return c.SendString("Login") })
	auth.Post("/refresh-token", func(c *fiber.Ctx) error { return c.SendString("Refresh JWT") })
	auth.Get("/me", func(c *fiber.Ctx) error { return c.SendString("Get My Profile & Role") })
	auth.Put("/profile", func(c *fiber.Ctx) error { return c.SendString("Update My Info") })
	auth.Put("/password", func(c *fiber.Ctx) error { return c.SendString("Update Password") })
	auth.Post("/forgot-password", func(c *fiber.Ctx) error { return c.SendString("Request Reset Link") })
	auth.Post("/reset-password", func(c *fiber.Ctx) error { return c.SendString("Reset Password with Token") })
	auth.Post("/verify-email", func(c *fiber.Ctx) error { return c.SendString("Verify Email") })
	auth.Post("/2fa/setup", func(c *fiber.Ctx) error { return c.SendString("Setup 2FA") })
	auth.Post("/logout", func(c *fiber.Ctx) error { return c.SendString("Logout Session") })

	// =========================================================================
	// 2. TENANT & COMPANY MANAGEMENT (8 API) - Tabel: Tenant, CompanySetting
	// =========================================================================
	tenant := v0.Group("/tenant")
	tenant.Get("/info", func(c *fiber.Ctx) error { return c.SendString("Get Company Info") })
	tenant.Put("/info", func(c *fiber.Ctx) error { return c.SendString("Update Company Info") })
	tenant.Get("/settings", func(c *fiber.Ctx) error { return c.SendString("Get Company Settings") })
	tenant.Put("/settings", func(c *fiber.Ctx) error { return c.SendString("Update Company Settings") })
	tenant.Get("/subscription", func(c *fiber.Ctx) error { return c.SendString("Get Current Plan") })
	tenant.Post("/subscription/upgrade", func(c *fiber.Ctx) error { return c.SendString("Upgrade Plan") })
	tenant.Get("/usage", func(c *fiber.Ctx) error { return c.SendString("Get OCR Usage Stats") })
	tenant.Delete("/terminate", func(c *fiber.Ctx) error { return c.SendString("Request Data Deletion") })

	// =========================================================================
	// 3. ORGANIZATION & USERS (15 API) - Tabel: Department, User, UserApprover
	// =========================================================================
	org := v0.Group("/org")
	// Departments
	org.Get("/departments", func(c *fiber.Ctx) error { return c.SendString("List Depts") })
	org.Post("/departments", func(c *fiber.Ctx) error { return c.SendString("Create Dept") })
	org.Get("/departments/:id", func(c *fiber.Ctx) error { return c.SendString("Detail Dept") })
	org.Put("/departments/:id", func(c *fiber.Ctx) error { return c.SendString("Update Dept") })
	org.Delete("/departments/:id", func(c *fiber.Ctx) error { return c.SendString("Delete Dept") })
	// Users
	org.Get("/users", func(c *fiber.Ctx) error { return c.SendString("List All Users") })
	org.Post("/users", func(c *fiber.Ctx) error { return c.SendString("Invite User via Email") })
	org.Get("/users/:id", func(c *fiber.Ctx) error { return c.SendString("Detail User") })
	org.Put("/users/:id", func(c *fiber.Ctx) error { return c.SendString("Update User Role/Dept") })
	org.Put("/users/:id/status", func(c *fiber.Ctx) error { return c.SendString("Activate/Deactivate User") })
	org.Delete("/users/:id", func(c *fiber.Ctx) error { return c.SendString("Soft Delete User") })
	// Approver & Hierarchy (UserApprover)
	org.Get("/hierarchy", func(c *fiber.Ctx) error { return c.SendString("Get Org Structure") })
	org.Get("/users/:id/approvers", func(c *fiber.Ctx) error { return c.SendString("Get User's Managers") })
	org.Post("/users/:id/approvers", func(c *fiber.Ctx) error { return c.SendString("Set User's Manager") })
	org.Delete("/users/:id/approvers/:approverId", func(c *fiber.Ctx) error { return c.SendString("Remove Approver") })

	// =========================================================================
	// 4. FINANCE MASTER DATA (16 API) - Tabel: AccountCategory, TaxRate, PaymentMethod, VendorMaster
	// =========================================================================
	finance := v0.Group("/finance")
	// Account Categories (Kanjo Kamoku)
	finance.Get("/categories", func(c *fiber.Ctx) error { return c.SendString("List Categories") })
	finance.Post("/categories", func(c *fiber.Ctx) error { return c.SendString("Create Category") })
	finance.Put("/categories/:id", func(c *fiber.Ctx) error { return c.SendString("Update Category") })
	finance.Delete("/categories/:id", func(c *fiber.Ctx) error { return c.SendString("Delete Category") })
	finance.Post("/categories/import", func(c *fiber.Ctx) error { return c.SendString("Bulk Import Categories") })
	// Tax Rates (Standard 10% vs Reduced 8%)
	finance.Get("/tax-rates", func(c *fiber.Ctx) error { return c.SendString("List Tax Rates") })
	finance.Post("/tax-rates", func(c *fiber.Ctx) error { return c.SendString("Create Tax Rate") })
	finance.Put("/tax-rates/:id", func(c *fiber.Ctx) error { return c.SendString("Update Tax Rate") })
	// Payment Methods
	finance.Get("/payments", func(c *fiber.Ctx) error { return c.SendString("List Payment Methods") })
	finance.Post("/payments", func(c *fiber.Ctx) error { return c.SendString("Create Payment Method") })
	// Vendor Master (Invois Seido Ready)
	finance.Get("/vendors", func(c *fiber.Ctx) error { return c.SendString("List Vendors") })
	finance.Post("/vendors", func(c *fiber.Ctx) error { return c.SendString("Create Vendor Master") })
	finance.Get("/vendors/:id", func(c *fiber.Ctx) error { return c.SendString("Detail Vendor") })
	finance.Put("/vendors/:id", func(c *fiber.Ctx) error { return c.SendString("Update Vendor Info") })
	finance.Delete("/vendors/:id", func(c *fiber.Ctx) error { return c.SendString("Delete Vendor") })
	finance.Get("/vendors/verify/:t_number", func(c *fiber.Ctx) error { return c.SendString("Verify T-Number via API Gov") })

	// =========================================================================
	// 5. RECEIPTS & AI-OCR (15 API) - Tabel: Receipt, ReceiptItem
	// =========================================================================
	receipts := v0.Group("/receipts")
	receipts.Get("/", handler.GetReceipts) // List Struk Pribadi
	receipts.Get("/all", func(c *fiber.Ctx) error { return c.SendString("List All (Admin Only)") })
	receipts.Post("/upload", func(c *fiber.Ctx) error { return c.SendString("Upload Image to OCR") })
	receipts.Post("/webhook/ocr", func(c *fiber.Ctx) error { return c.SendString("n8n/OpenAI Callback") })
	receipts.Get("/:id", func(c *fiber.Ctx) error { return c.SendString("Detail Struk & Item Detail") })
	receipts.Put("/:id", func(c *fiber.Ctx) error { return c.SendString("Confirm/Update OCR Data") })
	receipts.Delete("/:id", func(c *fiber.Ctx) error { return c.SendString("Soft Delete Struk") })
	receipts.Post("/bulk/delete", func(c *fiber.Ctx) error { return c.SendString("Bulk Delete Receipts") })
	receipts.Post("/bulk/update-category", func(c *fiber.Ctx) error { return c.SendString("Bulk Update Category") })
	receipts.Get("/:id/image", func(c *fiber.Ctx) error { return c.SendString("Get Secure Image URL") })
	receipts.Post("/:id/re-ocr", func(c *fiber.Ctx) error { return c.SendString("Re-trigger OCR Process") })
	receipts.Get("/search/advanced", func(c *fiber.Ctx) error { return c.SendString("Advanced Search Struk") })
	// Receipt Items
	receipts.Post("/:id/items", func(c *fiber.Ctx) error { return c.SendString("Add Manual Item") })
	receipts.Put("/items/:itemId", func(c *fiber.Ctx) error { return c.SendString("Update Item Detail") })
	receipts.Delete("/items/:itemId", func(c *fiber.Ctx) error { return c.SendString("Delete Item") })

	// =========================================================================
	// 6. EXPENSE REPORTS (10 API) - Tabel: ExpenseReport
	// =========================================================================
	reports := v0.Group("/reports")
	reports.Get("/", func(c *fiber.Ctx) error { return c.SendString("List My Reports") })
	reports.Post("/", func(c *fiber.Ctx) error { return c.SendString("Create Bundle Report") })
	reports.Get("/:id", func(c *fiber.Ctx) error { return c.SendString("Detail Report") })
	reports.Put("/:id", func(c *fiber.Ctx) error { return c.SendString("Update Report Title/Data") })
	reports.Delete("/:id", func(c *fiber.Ctx) error { return c.SendString("Delete Report") })
	reports.Post("/:id/submit", func(c *fiber.Ctx) error { return c.SendString("Submit to Workflow Approval") })
	reports.Post("/:id/cancel", func(c *fiber.Ctx) error { return c.SendString("Cancel Submitted Report") })
	reports.Get("/stats/monthly", func(c *fiber.Ctx) error { return c.SendString("Monthly Report Stats") })
	reports.Get("/export/preview", func(c *fiber.Ctx) error { return c.SendString("Preview PDF of Report") })
	reports.Post("/bulk/submit", func(c *fiber.Ctx) error { return c.SendString("Bulk Submit Reports") })

	// =========================================================================
	// 7. WORKFLOW & APPROVALS (12 API) - Tabel: ApprovalWorkflow, Step, Log
	// =========================================================================
	approvals := v0.Group("/approvals")
	approvals.Get("/pending", func(c *fiber.Ctx) error { return c.SendString("List Tasks for Me") })
	approvals.Post("/action", func(c *fiber.Ctx) error { return c.SendString("Approve/Reject Action") })
	approvals.Post("/remand", func(c *fiber.Ctx) error { return c.SendString("Return to Employee (Remand)") })
	approvals.Get("/history", func(c *fiber.Ctx) error { return c.SendString("My Approval History") })
	// Config Workflow (Admin)
	approvals.Get("/workflows", func(c *fiber.Ctx) error { return c.SendString("List Workflow Templates") })
	approvals.Post("/workflows", func(c *fiber.Ctx) error { return c.SendString("Create New Workflow Template") })
	approvals.Put("/workflows/:id", func(c *fiber.Ctx) error { return c.SendString("Update Template") })
	approvals.Get("/workflows/:id/steps", func(c *fiber.Ctx) error { return c.SendString("Get Steps Detail") })
	approvals.Post("/workflows/:id/steps", func(c *fiber.Ctx) error { return c.SendString("Add Step Approver") })
	approvals.Put("/steps/:stepId", func(c *fiber.Ctx) error { return c.SendString("Update Step Order/Approver") })
	approvals.Delete("/steps/:stepId", func(c *fiber.Ctx) error { return c.SendString("Remove Step") })
	approvals.Get("/delegations", func(c *fiber.Ctx) error { return c.SendString("Get Proxy Approver (Delegation)") })

	// =========================================================================
	// 8. COMPLIANCE & AUDIT (10 API) - Tabel: AuditTrail, ExportLog
	// =========================================================================
	audit := v0.Group("/audit")
	audit.Get("/summary", handler.GetStats) // Dashboard Summary
	audit.Get("/trails", func(c *fiber.Ctx) error { return c.SendString("List All Audit Trails") })
	audit.Get("/trails/receipt/:id", func(c *fiber.Ctx) error { return c.SendString("Get History of Specific Receipt") })
	audit.Get("/trails/user/:id", func(c *fiber.Ctx) error { return c.SendString("Get History by User") })
	// Exports
	audit.Post("/export/csv", func(c *fiber.Ctx) error { return c.SendString("Generate CSV for freee/MoneyForward") })
	audit.Post("/export/pdf-bundle", func(c *fiber.Ctx) error { return c.SendString("Generate Audit PDF Bundle") })
	audit.Get("/export/logs", func(c *fiber.Ctx) error { return c.SendString("List Export History") })
	audit.Get("/export/download/:id", func(c *fiber.Ctx) error { return c.SendString("Secure Download Link") })
	audit.Get("/tax-summary", func(c *fiber.Ctx) error { return c.SendString("Tax 8% vs 10% Summary Report") })
	audit.Post("/data-integrity/check", func(c *fiber.Ctx) error { return c.SendString("Check Data Hash Integrity") })

	// =========================================================================
	// 9. SYSTEM ADMINISTRATION (Super Admin) (10 API) - Tabel: SubscriptionPlan
	// =========================================================================
	sys := v0.Group("/system")
	sys.Get("/plans", func(c *fiber.Ctx) error { return c.SendString("List Plans (Free/Pro)") })
	sys.Post("/plans", func(c *fiber.Ctx) error { return c.SendString("Create Subscription Plan") })
	sys.Put("/plans/:id", func(c *fiber.Ctx) error { return c.SendString("Update Plan Detail") })
	sys.Get("/tenants", func(c *fiber.Ctx) error { return c.SendString("List All Tenants Global") })
	sys.Put("/tenants/:id/status", func(c *fiber.Ctx) error { return c.SendString("Activate/Suspend Tenant") })
	sys.Get("/monitoring/ocr", func(c *fiber.Ctx) error { return c.SendString("Monitor AI-OCR Health") })
	sys.Get("/monitoring/errors", func(c *fiber.Ctx) error { return c.SendString("System Error Logs") })
	sys.Post("/maintenance/mode", func(c *fiber.Ctx) error { return c.SendString("Toggle Maintenance Mode") })
	sys.Get("/backup/configs", func(c *fiber.Ctx) error { return c.SendString("Backup Configurations") })
	sys.Get("/version", func(c *fiber.Ctx) error { return c.SendString("Get System Version") })
}
*/
