package configs

import (
	"fmt"
	"log"
	"os"
	"time"

	"ocr-saas-backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB(cfg *Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode, cfg.DBTimeZone,
	)

	// --- KONFIGURASI LOGGER CUSTOM ---
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // 👈 Naikkan ke 1 detik biar terminal gak berisik
			LogLevel:                  logger.Warn, // Munculkan hanya Warning & Error
			IgnoreRecordNotFoundError: true,        // Abaikan log "record not found"
			Colorful:                  true,        // Pakai warna di terminal
		},
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger, // Gunakan custom logger
	})

	if err != nil {
		log.Fatalf("❌ DB Connection Failed: %v", err)
	}

	sqlDB, _ := DB.DB()
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	fmt.Println("🚀 Syncing Database Models...")
	err = DB.AutoMigrate(
		&models.SubscriptionPlan{}, &models.Tenant{}, &models.CompanySetting{},
		&models.Department{}, &models.User{}, &models.RefreshToken{},
		&models.UserApprover{}, &models.AccountCategory{}, &models.TaxRate{},
		&models.PaymentMethod{}, &models.VendorMaster{}, &models.ExpenseReport{},
		&models.ApprovalWorkflow{}, &models.ApprovalStep{}, &models.Receipt{},
		&models.ReceiptItem{}, &models.ApprovalLog{}, &models.AuditTrail{},
		&models.ExportLog{}, &models.TenantUsage{}, &models.OCRJob{},
		&models.Budget{},
	)

	if err != nil {
		log.Fatalf("❌ Migration Failed: %v", err)
	}

	log.Println("✅ DB Connected & Migrated")
	return DB
}
