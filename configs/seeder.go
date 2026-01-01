package configs

import (
	"fmt"
	"ocr-saas-backend/internal/models"
	"time"

	"golang.org/x/crypto/bcrypt" // WAJIB ADA
	"gorm.io/gorm"
)

func SeedDatabase(db *gorm.DB) {
	var count int64
	// Cek apakah database sudah ada datanya (berdasarkan Tenant)
	db.Model(&models.Tenant{}).Count(&count)
	if count > 0 {
		fmt.Println("✅ Database sudah terisi, skip seeding.")
		return
	}

	fmt.Println("🚀 Memulai Seeding 18 Tabel Enterprise...")

	// --- 1. SUBSCRIPTION PLAN ---
	plan := models.SubscriptionPlan{
		Name:        "Business Enterprise",
		MaxReceipts: 5000,
		Price:       15000,
	}
	db.Create(&plan)

	// --- 2. TENANT (PERUSAHAAN) ---
	tenant := models.Tenant{
		Name:               "株式会社テクノロジー・ジャパン (Tech Japan Corp)",
		Subdomain:          "tech-japan",
		SubscriptionPlanID: plan.ID,
		BusinessNumber:     "5011001043210", // Hojin Bango (13 digit)
		Status:             "ACTIVE",
	}
	db.Create(&tenant)

	// --- 3. COMPANY SETTINGS ---
	db.Create(&models.CompanySetting{
		TenantID:   tenant.ID,
		DateFormat: "YYYY/MM/DD",
		Currency:   "JPY",
		AutoOCR:    true,
	})

	// --- 4. DEPARTMENTS ---
	salesDept := models.Department{
		TenantID: tenant.ID,
		Name:     "営業部 (Sales Department)",
		Code:     "SLS01",
	}
	db.Create(&salesDept)

	financeDept := models.Department{
		TenantID: tenant.ID,
		Name:     "経理部 (Finance Department)",
		Code:     "FIN01",
	}
	db.Create(&financeDept)

	// Hash Password untuk login
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), 10)
	passString := string(hashedPassword)

	// --- 5. USERS (Admin, Manager, Staff) ---
	admin := models.User{
		TenantID:     tenant.ID,
		DepartmentID: &financeDept.ID,
		Name:         "佐藤 健一 (Sato Admin)",
		Email:        "admin@tech-japan.jp",
		Role:         "ADMIN",
		PasswordHash: passString,
	}
	db.Create(&admin)

	manager := models.User{
		TenantID:     tenant.ID,
		DepartmentID: &salesDept.ID,
		Name:         "鈴木 一郎 (Suzuki Manager)",
		Email:        "manager@tech-japan.jp",
		Role:         "MANAGER",
		PasswordHash: passString,
	}
	db.Create(&manager)

	staff := models.User{
		TenantID:     tenant.ID,
		DepartmentID: &salesDept.ID,
		Name:         "田中 太郎 (Tanaka Staff)",
		Email:        "staff@tech-japan.jp",
		Role:         "EMPLOYEE",
		PasswordHash: passString,
	}
	db.Create(&staff) // TADI BAGIAN INI HILANG

	// --- 6. USER APPROVER (Relasi Bawahan -> Atasan) ---
	db.Create(&models.UserApprover{
		EmployeeID: staff.ID,
		ApproverID: manager.ID,
	})

	// --- 7. ACCOUNT CATEGORIES ---
	travelCat := models.AccountCategory{
		TenantID: tenant.ID,
		Code:     "101",
		Name:     "旅費交通費 (Travel Expense)",
	}
	db.Create(&travelCat)

	mealCat := models.AccountCategory{
		TenantID: tenant.ID,
		Code:     "102",
		Name:     "会議費/食事代 (Meal/Meeting)",
	}
	db.Create(&mealCat)

	// --- 8. TAX RATES ---
	taxStandard := models.TaxRate{TenantID: tenant.ID, Name: "Standard 10%", Percentage: 10}
	taxReduced := models.TaxRate{TenantID: tenant.ID, Name: "Reduced 8%", Percentage: 8}
	db.Create(&taxStandard)
	db.Create(&taxReduced)

	// --- 9. PAYMENT METHODS ---
	db.Create(&models.PaymentMethod{TenantID: tenant.ID, Name: "Corporate Card (Visa)"})
	db.Create(&models.PaymentMethod{TenantID: tenant.ID, Name: "Cash (現金)"})

	// --- 10. VENDOR MASTER ---
	vendor := models.VendorMaster{
		TenantID:  tenant.ID,
		Name:      "Lawson Shibuya Station",
		TaxNumber: "T1234567890123",
	}
	db.Create(&vendor)

	// --- 11. APPROVAL WORKFLOW ---
	workflow := models.ApprovalWorkflow{
		TenantID: tenant.ID,
		Name:     "Standard Expense Workflow",
		IsActive: true,
	}
	db.Create(&workflow)

	// --- 12. APPROVAL STEPS ---
	db.Create(&models.ApprovalStep{
		WorkflowID: workflow.ID,
		StepOrder:  1,
		ApproverID: manager.ID,
	})

	// --- 13. EXPENSE REPORT ---
	report := models.ExpenseReport{
		TenantID:    tenant.ID,
		UserID:      staff.ID,
		Title:       "出張費用精算 - 2025年10月",
		TotalAmount: 2500,
		Status:      "PENDING",
	}
	db.Create(&report)

	// --- 14. RECEIPT ---
	now := time.Now()
	receipt := models.Receipt{
		TenantID:          tenant.ID,
		UserID:            staff.ID,
		ReportID:          &report.ID,
		AccountCategoryID: &travelCat.ID,
		StoreName:         "ファミリーマート (FamilyMart)",
		TransactionDate:   &now,
		TotalAmount:       1500,
		TaxRegistrationID: "T5011001043210",
		IsQualified:       true,
		Status:            "PENDING",
		ImageURL:          "https://storage.googleapis.com/demo/receipt_001.jpg",
	}
	db.Create(&receipt)

	// --- 15. RECEIPT ITEMS ---
	db.Create(&models.ReceiptItem{
		ReceiptID:   receipt.ID,
		Description: "新幹線チケット",
		Amount:      1364,
		TaxAmount:   136,
		TaxRate:     10,
	})

	// --- 16. APPROVAL LOG ---
	db.Create(&models.ApprovalLog{
		ExpenseReportID: &report.ID,
		UserID:          staff.ID,
		Action:          "SUBMIT",
		Comment:         "大阪出張の精算をお願いいたします。",
	})

	// --- 17. AUDIT TRAIL ---
	db.Create(&models.AuditTrail{
		TenantID:  tenant.ID,
		UserID:    staff.ID,
		Action:    "CREATE_REPORT",
		TableName: "expense_reports",
		RecordID:  report.ID.String(),
		OldData:   "{}",
		NewData:   "{\"title\":\"出張費用精算\"}",
	})

	// --- 18. EXPORT LOG ---
	db.Create(&models.ExportLog{
		TenantID: tenant.ID,
		UserID:   admin.ID,
		Format:   "CSV (Money Forward Format)",
		FileURL:  "https://storage.googleapis.com/demo/export_oct.csv",
	})

	fmt.Println("✅ Seeding Selesai! 18 Tabel siap digunakan untuk Demo.")
}