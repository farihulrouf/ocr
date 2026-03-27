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

var DB *gorm.DB

func ConnectDB(cfg *Config) *gorm.DB {
	fmt.Println("🔥 TRY CONNECT DB...")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
		cfg.DBSSLMode,
		cfg.DBTimeZone,
	)

	// 🔥 SAFE DEBUG
	fmt.Println("DSN DEBUG:",
		fmt.Sprintf("host=%s user=%s dbname=%s port=%s sslmode=%s",
			cfg.DBHost,
			cfg.DBUser,
			cfg.DBName,
			cfg.DBPort,
			cfg.DBSSLMode,
		),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("❌ Gagal konek DB: %v", err)
	}

	log.Println("✅ GORM connected")

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("❌ Gagal ambil sql.DB: %v", err)
	}

	// 🔥 WAJIB — biar gak false positive
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("❌ DB PING FAILED: %v", err)
	}

	log.Println("✅ DB PING SUCCESS")

	// Pool config
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
