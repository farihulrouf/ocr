package configs

import (
	"fmt"
	"ocr-saas-backend/internal/models"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

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

	vendor1 := models.VendorMaster{TenantID: tenant1.ID, Name: "Lawson Shibuya", TaxNumber: "T1234567890123"}
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
		TenantID: tenant1.ID, UserID: staff1.ID,
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

	db.Create(&models.ApprovalLog{
		ExpenseReportID: &rep1.ID,
		UserID:          staff1.ID,
		Action:          "SUBMIT",
		Comment:         "お願いします。",
	})

	db.Create(&models.AuditTrail{
		TenantID:  tenant1.ID,
		UserID:    staff1.ID,
		Action:    "CREATE_REPORT",
		TableName: "expense_reports",
		RecordID:  rep1.ID.String(),
		OldData:   "{}",
		NewData:   "{\"title\":\"出張費用精算\"}",
	})

	db.Create(&models.ExportLog{
		TenantID: tenant1.ID,
		UserID:   admin1.ID,
		Format:   "CSV",
		FileURL:  "https://storage.googleapis.com/demo/t1_export.csv",
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

	db.Create(&models.UserApprover{EmployeeID: staff2.ID, ApproverID: manager2.ID})

	cat2A := models.AccountCategory{TenantID: tenant2.ID, Code: "201", Name: "仕入れ費用"}
	cat2B := models.AccountCategory{TenantID: tenant2.ID, Code: "202", Name: "調理用品費"}
	db.Create(&cat2A)
	db.Create(&cat2B)

	tax2A := models.TaxRate{TenantID: tenant2.ID, Name: "10%", Percentage: 10}
	tax2B := models.TaxRate{TenantID: tenant2.ID, Name: "軽減8%", Percentage: 8}
	db.Create(&tax2A)
	db.Create(&tax2B)

	db.Create(&models.PaymentMethod{TenantID: tenant2.ID, Name: "Bank Transfer"})
	db.Create(&models.PaymentMethod{TenantID: tenant2.ID, Name: "Cash"})

	vendor2 := models.VendorMaster{TenantID: tenant2.ID, Name: "Aeon Market", TaxNumber: "T9876543210001"}
	db.Create(&vendor2)

	wf2 := models.ApprovalWorkflow{TenantID: tenant2.ID, Name: "Food Approval Flow", IsActive: true}
	db.Create(&wf2)
	db.Create(&models.ApprovalStep{WorkflowID: wf2.ID, StepOrder: 1, ApproverID: manager2.ID})

	rep2 := models.ExpenseReport{
		TenantID: tenant2.ID, UserID: staff2.ID,
		Title:       "食材仕入れ精算",
		TotalAmount: 9000, Status: "PENDING",
	}
	db.Create(&rep2)

	rc2 := models.Receipt{
		TenantID: tenant2.ID, UserID: staff2.ID,
		ReportID:          &rep2.ID,
		AccountCategoryID: &cat2A.ID,
		StoreName:         "イオンマーケット",
		TransactionDate:   &now,
		TotalAmount:       8700,
		TaxRegistrationID: "T7012002045678",
		IsQualified:       true,
		Status:            "PENDING",
		ImageURL:          "https://storage.googleapis.com/demo/r2.jpg",
	}
	db.Create(&rc2)

	db.Create(&models.ReceiptItem{
		ReceiptID:   rc2.ID,
		Description: "野菜・肉購入",
		Amount:      8000,
		TaxAmount:   700,
		TaxRate:     10,
	})

	db.Create(&models.ApprovalLog{
		ExpenseReportID: &rep2.ID,
		UserID:          staff2.ID,
		Action:          "SUBMIT",
		Comment:         "本日の仕入れです。",
	})

	db.Create(&models.AuditTrail{
		TenantID:  tenant2.ID,
		UserID:    staff2.ID,
		Action:    "CREATE_REPORT",
		TableName: "expense_reports",
		RecordID:  rep2.ID.String(),
		OldData:   "{}",
		NewData:   "{\"title\":\"仕入れ\"}",
	})

	db.Create(&models.ExportLog{
		TenantID: tenant2.ID,
		UserID:   admin2.ID,
		Format:   "CSV",
		FileURL:  "https://storage.googleapis.com/demo/t2_export.csv",
	})

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

	db.Create(&models.CompanySetting{
		TenantID: tenant3.ID, Currency: "JPY", DateFormat: "DD-MM-YYYY", AutoOCR: true,
	})

	dept3A := models.Department{TenantID: tenant3.ID, Name: "設計部", Code: "DSN01"}
	dept3B := models.Department{TenantID: tenant3.ID, Name: "管理部", Code: "ADM01"}
	db.Create(&dept3A)
	db.Create(&dept3B)

	admin3 := models.User{
		TenantID: tenant3.ID, DepartmentID: &dept3B.ID,
		Name: "高橋 良介", Email: "admin@kyoto-eng.jp", Role: "ADMIN", PasswordHash: pass,
	}
	manager3 := models.User{
		TenantID: tenant3.ID, DepartmentID: &dept3A.ID,
		Name: "藤田 直樹", Email: "manager@kyoto-eng.jp", Role: "MANAGER", PasswordHash: pass,
	}
	staff3 := models.User{
		TenantID: tenant3.ID, DepartmentID: &dept3A.ID,
		Name: "松本 明", Email: "staff@kyoto-eng.jp", Role: "EMPLOYEE", PasswordHash: pass,
	}
	db.Create(&admin3)
	db.Create(&manager3)
	db.Create(&staff3)

	db.Create(&models.UserApprover{EmployeeID: staff3.ID, ApproverID: manager3.ID})

	cat3A := models.AccountCategory{TenantID: tenant3.ID, Code: "301", Name: "研究開発費"}
	cat3B := models.AccountCategory{TenantID: tenant3.ID, Code: "302", Name: "備品費用"}
	db.Create(&cat3A)
	db.Create(&cat3B)

	tax3A := models.TaxRate{TenantID: tenant3.ID, Name: "10%", Percentage: 10}
	tax3B := models.TaxRate{TenantID: tenant3.ID, Name: "非課税", Percentage: 0}
	db.Create(&tax3A)
	db.Create(&tax3B)

	db.Create(&models.PaymentMethod{TenantID: tenant3.ID, Name: "Cash"})
	db.Create(&models.PaymentMethod{TenantID: tenant3.ID, Name: "Corporate Card"})

	vendor3 := models.VendorMaster{TenantID: tenant3.ID, Name: "Yodobashi Kyoto", TaxNumber: "T5555555554321"}
	db.Create(&vendor3)

	wf3 := models.ApprovalWorkflow{TenantID: tenant3.ID, Name: "Engineering Flow", IsActive: true}
	db.Create(&wf3)
	db.Create(&models.ApprovalStep{WorkflowID: wf3.ID, StepOrder: 1, ApproverID: manager3.ID})

	rep3 := models.ExpenseReport{
		TenantID: tenant3.ID, UserID: staff3.ID,
		Title:       "機材購入精算",
		TotalAmount: 12000, Status: "PENDING",
	}
	db.Create(&rep3)

	rc3 := models.Receipt{
		TenantID: tenant3.ID, UserID: staff3.ID,
		ReportID: &rep3.ID, AccountCategoryID: &cat3B.ID,
		StoreName:         "ヨドバシカメラ 京都",
		TransactionDate:   &now,
		TotalAmount:       11800,
		TaxRegistrationID: "T3015003098765",
		IsQualified:       true,
		Status:            "PENDING",
		ImageURL:          "https://storage.googleapis.com/demo/r3.jpg",
	}
	db.Create(&rc3)

	db.Create(&models.ReceiptItem{
		ReceiptID:   rc3.ID,
		Description: "精密ドライバーセット",
		Amount:      11800,
		TaxAmount:   0,
		TaxRate:     0,
	})

	db.Create(&models.ApprovalLog{
		ExpenseReportID: &rep3.ID,
		UserID:          staff3.ID,
		Action:          "SUBMIT",
		Comment:         "備品です。",
	})

	db.Create(&models.AuditTrail{
		TenantID:  tenant3.ID,
		UserID:    staff3.ID,
		Action:    "CREATE_REPORT",
		TableName: "expense_reports",
		RecordID:  rep3.ID.String(),
		OldData:   "{}",
		NewData:   "{\"title\":\"備品購入\"}",
	})

	db.Create(&models.ExportLog{
		TenantID: tenant3.ID,
		UserID:   admin3.ID,
		Format:   "CSV",
		FileURL:  "https://storage.googleapis.com/demo/t3_export.csv",
	})

	fmt.Println("🎉 Selesai! 3 Tenant & 18 Tabel berhasil di-seed.")
}
