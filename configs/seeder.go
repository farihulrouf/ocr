package configs

import (
	"fmt"
	"ocr-saas-backend/internal/models"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Helper function untuk mengubah UUID menjadi *uuid.UUID (Pointer)
// Karena field TenantID di VendorMaster sekarang bisa NULL (Pointer)
func uuidPtr(id uuid.UUID) *uuid.UUID {
	return &id
}

func SeedDatabase(db *gorm.DB) {
	var count int64
	db.Model(&models.Tenant{}).Count(&count)
	if count > 0 {
		fmt.Println("✅ Database sudah terisi, skip seeding.")
		return
	}

	fmt.Println("🚀 Memulai Seeding 3 Tenant · 18 Tabel Lengkap")

	// HASH PASSWORD SEKALI SAJA
	hashed, _ := bcrypt.GenerateFromPassword([]byte("password123"), 10)
	pass := string(hashed)

	// ============================================================
	// ================== TENANT 1: TECH JAPAN =====================
	// ============================================================

	plan := models.SubscriptionPlan{
		Name:        "Business Enterprise",
		MaxReceipts: 5000,
		Price:       15000,
	}
	db.Create(&plan)

	tenant1 := models.Tenant{
		Name:               "株式会社テクノロジー・ジャパン (Tech Japan Corp)",
		Subdomain:          "tech-japan",
		SubscriptionPlanID: plan.ID,
		BusinessNumber:     "5011001043210",
		Status:             "ACTIVE",
	}
	db.Create(&tenant1)

	db.Create(&models.CompanySetting{
		TenantID:   tenant1.ID,
		DateFormat: "YYYY/MM/DD",
		Currency:   "JPY",
		AutoOCR:    true,
	})

	dept1A := models.Department{TenantID: tenant1.ID, Name: "営業部", Code: "SLS01"}
	dept1B := models.Department{TenantID: tenant1.ID, Name: "経理部", Code: "FIN01"}
	db.Create(&dept1A)
	db.Create(&dept1B)

	admin1 := models.User{
		TenantID: tenant1.ID, DepartmentID: &dept1B.ID,
		Name: "佐藤 健一", Email: "admin@tech-japan.jp",
		Role: "ADMIN", PasswordHash: pass,
	}
	manager1 := models.User{
		TenantID: tenant1.ID, DepartmentID: &dept1A.ID,
		Name: "鈴木 一郎", Email: "manager@tech-japan.jp",
		Role: "MANAGER", PasswordHash: pass,
	}
	staff1 := models.User{
		TenantID: tenant1.ID, DepartmentID: &dept1A.ID,
		Name: "田中 太郎", Email: "staff@tech-japan.jp",
		Role: "EMPLOYEE", PasswordHash: pass,
	}
	db.Create(&admin1)
	db.Create(&manager1)
	db.Create(&staff1)

	db.Create(&models.UserApprover{EmployeeID: staff1.ID, ApproverID: manager1.ID})

	cat1A := models.AccountCategory{TenantID: tenant1.ID, Code: "101", Name: "旅費交通費"}
	cat1B := models.AccountCategory{TenantID: tenant1.ID, Code: "102", Name: "会議費"}
	db.Create(&cat1A)
	db.Create(&cat1B)

	tax1A := models.TaxRate{TenantID: tenant1.ID, Name: "10%", Percentage: 10}
	tax1B := models.TaxRate{TenantID: tenant1.ID, Name: "8%", Percentage: 8}
	db.Create(&tax1A)
	db.Create(&tax1B)

	db.Create(&models.PaymentMethod{TenantID: tenant1.ID, Name: "Visa Corporate Card"})
	db.Create(&models.PaymentMethod{TenantID: tenant1.ID, Name: "Cash"})

	// FIX: Pakai uuidPtr
	vendor1 := models.VendorMaster{
		TenantID:  uuidPtr(tenant1.ID),
		Name:      "Lawson Shibuya",
		TaxNumber: "T1234567890123",
	}
	db.Create(&vendor1)

	wf1 := models.ApprovalWorkflow{TenantID: tenant1.ID, Name: "Standard Workflow", IsActive: true}
	db.Create(&wf1)
	db.Create(&models.ApprovalStep{WorkflowID: wf1.ID, StepOrder: 1, ApproverID: manager1.ID})

	rep1 := models.ExpenseReport{
		TenantID: tenant1.ID, UserID: staff1.ID,
		Title:       "出張費用精算 (大阪)",
		TotalAmount: 2500, Status: "PENDING",
	}
	db.Create(&rep1)

	now := time.Now()
	rc1 := models.Receipt{
		TenantID:          tenant1.ID,
		UserID:            staff1.ID,
		ReportID:          &rep1.ID,
		AccountCategoryID: &cat1A.ID,
		StoreName:         "ファミリーマート",
		TransactionDate:   &now,
		TotalAmount:       1500,
		TaxRegistrationID: "T5011001043210",
		IsQualified:       true,
		Status:            "PENDING",
		ImageURL:          "https://storage.googleapis.com/demo/r1.jpg",
	}
	db.Create(&rc1)

	db.Create(&models.ReceiptItem{
		ReceiptID:   rc1.ID,
		Description: "新幹線チケット",
		Amount:      1364,
		TaxAmount:   136,
		TaxRate:     10,
	})

	// ============================================================
	// ================== TENANT 2: TOKYO FOOD ====================
	// ============================================================

	tenant2 := models.Tenant{
		Name:               "東京フードサービス株式会社 (Tokyo Food Service)",
		Subdomain:          "tokyo-food",
		SubscriptionPlanID: plan.ID,
		BusinessNumber:     "7012002045678",
		Status:             "ACTIVE",
	}
	db.Create(&tenant2)

	db.Create(&models.CompanySetting{
		TenantID: tenant2.ID, Currency: "JPY", DateFormat: "YYYY-MM-DD", AutoOCR: false,
	})

	dept2A := models.Department{TenantID: tenant2.ID, Name: "店舗管理部", Code: "STO01"}
	dept2B := models.Department{TenantID: tenant2.ID, Name: "会計部", Code: "ACC01"}
	db.Create(&dept2A)
	db.Create(&dept2B)

	admin2 := models.User{
		TenantID: tenant2.ID, DepartmentID: &dept2B.ID,
		Name: "山田 花子", Email: "admin@tokyo-food.jp", Role: "ADMIN", PasswordHash: pass,
	}
	manager2 := models.User{
		TenantID: tenant2.ID, DepartmentID: &dept2A.ID,
		Name: "中村 大輔", Email: "manager@tokyo-food.jp", Role: "MANAGER", PasswordHash: pass,
	}
	staff2 := models.User{
		TenantID: tenant2.ID, DepartmentID: &dept2A.ID,
		Name: "吉田 美咲", Email: "staff@tokyo-food.jp", Role: "EMPLOYEE", PasswordHash: pass,
	}
	db.Create(&admin2)
	db.Create(&manager2)
	db.Create(&staff2)

	cat2A := models.AccountCategory{TenantID: tenant2.ID, Code: "201", Name: "仕入れ費用"}
	db.Create(&cat2A)

	// FIX: Pakai uuidPtr
	vendor2 := models.VendorMaster{
		TenantID:  uuidPtr(tenant2.ID),
		Name:      "Aeon Market",
		TaxNumber: "T9876543210001",
	}
	db.Create(&vendor2)

	// ============================================================
	// ================= TENANT 3: KYOTO ENGINEERING ==============
	// ============================================================

	tenant3 := models.Tenant{
		Name:               "京都エンジニアリング株式会社 (Kyoto Engineering)",
		Subdomain:          "kyoto-eng",
		SubscriptionPlanID: plan.ID,
		BusinessNumber:     "3015003098765",
		Status:             "ACTIVE",
	}
	db.Create(&tenant3)

	dept3A := models.Department{TenantID: tenant3.ID, Name: "設計部", Code: "DSN01"}
	db.Create(&dept3A)

	staff3 := models.User{
		TenantID: tenant3.ID, DepartmentID: &dept3A.ID,
		Name: "松本 明", Email: "staff@kyoto-eng.jp", Role: "EMPLOYEE", PasswordHash: pass,
	}
	db.Create(&staff3)

	// FIX: Pakai uuidPtr
	vendor3 := models.VendorMaster{
		TenantID:  uuidPtr(tenant3.ID),
		Name:      "Yodobashi Kyoto",
		TaxNumber: "T5555555554321",
	}
	db.Create(&vendor3)

	// ============================================================
	// =================== GLOBAL VENDOR SEEDING ==================
	// ============================================================
	// Vendor yang tidak punya TenantID (NULL) agar bisa dipakai semua
	fmt.Println("🌍 Seeding Global Vendors (Indomaret, Starbucks, etc)...")

	globalVendors := []models.VendorMaster{
		{
			Name:        "PT Sari Coffee Indonesia",
			DisplayName: "Starbucks",
			Aliases:     "SBUX, Starbucks Coffee, Starbucks Indonesia",
			TaxNumber:   "01.328.014.1-054.000",
			Category:    "Food & Beverage",
			CountryCode: "ID",
			IsGlobal:    true,
		},
		{
			Name:        "PT Indomarco Prismatama",
			DisplayName: "Indomaret",
			Aliases:     "Indomaret, IDM, Indomarco",
			TaxNumber:   "01.337.994.1-092.000",
			Category:    "Groceries",
			CountryCode: "ID",
			IsGlobal:    true,
		},
		{
			Name:        "PT Sumber Alfaria Trijaya Tbk",
			DisplayName: "Alfamart",
			Aliases:     "Alfamart, ALFA, SAT",
			TaxNumber:   "01.334.885.1-054.000",
			Category:    "Groceries",
			CountryCode: "ID",
			IsGlobal:    true,
		},
		{
			Name:        "PT Pertamina Patra Niaga",
			DisplayName: "Pertamina",
			Aliases:     "Pertamina, SPBU, Pertamax",
			TaxNumber:   "01.000.222.3-000.000",
			Category:    "Fuel",
			CountryCode: "ID",
			IsGlobal:    true,
		},
	}

	for _, gv := range globalVendors {
		db.Create(&gv)
	}

	// ============================================================
	// ===================== TAMBAHAN LOGS ========================
	// ============================================================

	db.Create(&models.TenantUsage{
		TenantID: tenant1.ID,
		OCRLimit: 5000,
		OCRUsed:  120,
	})

	// Refresh Tokens
	db.Create(&models.RefreshToken{
		ID:        uuid.New(),
		UserID:    admin1.ID,
		Token:     "refresh_token_admin_1",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})

	fmt.Println("🎉 Selesai! 3 Tenant & Global Vendors berhasil di-seed.")
}
