package configs

import (
	"fmt"
	"log"
	"time"

	"ocr-saas-backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// masih pakai global (nanti kita hapus di step berikutnya)
var DB *gorm.DB

func ConnectDB(cfg *Config) *gorm.DB {
	host := cfg.DBHost
	user := cfg.DBUser
	password := cfg.DBPassword
	dbname := cfg.DBName
	port := cfg.DBPort

	sslmode := "disable"

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		host, user, password, dbname, port, sslmode,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})

	if err != nil {
		log.Fatalf("❌ Gagal konek DB: %v", err)
	}

	log.Println("✅ Database connected")

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("❌ Gagal ambil sql.DB: %v", err)
	}

	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("✅ Connection pool ready")

	// 🔥 Auto migrate (DEV ONLY)
	err = DB.AutoMigrate(
		&models.SubscriptionPlan{},
		&models.Tenant{},
		&models.CompanySetting{},
		&models.Department{},
		&models.User{},
		&models.RefreshToken{},
		&models.UserApprover{},
		&models.AccountCategory{},
		&models.TaxRate{},
		&models.PaymentMethod{},
		&models.VendorMaster{},
		&models.ExpenseReport{},
		&models.ApprovalWorkflow{},
		&models.ApprovalStep{},
		&models.Receipt{},
		&models.ReceiptItem{},
		&models.ApprovalLog{},
		&models.AuditTrail{},
		&models.ExportLog{},
		&models.TenantUsage{},
		&models.OCRJob{},
	)

	if err != nil {
		log.Fatalf("❌ Gagal migrate: %v", err)
	}

	log.Println("✅ Migration done")

	return DB
}
