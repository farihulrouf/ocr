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

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	var err error
	// 🔄 RETRY LOGIC: Menunggu Postgres siap (penting untuk Docker)
	for i := 0; i < 5; i++ {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: newLogger,
		})
		if err == nil {
			break
		}
		log.Printf("⏳ Waiting for Database... (Attempt %d/5)", i+1)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Fatalf("❌ DB Connection Failed: %v", err)
	}

	sqlDB, _ := DB.DB()
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 🚩 KONTROL MIGRASI: Lewati migrasi jika sedang melakukan Restore manual
	// Kamu bisa menambahkan env variable SKIP_MIGRATION jika ingin mematikan ini total
	if os.Getenv("SKIP_MIGRATION") != "true" {
		fmt.Println("🚀 Syncing Database Models...")
		err = DB.AutoMigrate(
			&models.SubscriptionPlan{}, &models.Tenant{}, &models.CompanySetting{},
			&models.Department{}, &models.User{}, &models.RefreshToken{},
			&models.UserApprover{}, &models.AccountCategory{}, &models.TaxRate{},
			&models.PaymentMethod{}, &models.VendorMaster{}, &models.ExpenseReport{},
			&models.ApprovalWorkflow{}, &models.ApprovalStep{}, &models.Receipt{},
			&models.ReceiptItem{}, &models.ApprovalLog{}, &models.AuditTrail{},
			&models.ExportLog{}, &models.TenantUsage{}, &models.OCRJob{},
			&models.Budget{}, &models.Disbursement{},
		)

		if err != nil {
			// Jika error karena sudah ada (saat restore), kita log saja, jangan Fatal
			log.Printf("⚠️  Migration Warning: %v (This is normal if you are restoring data)", err)
		}
	}

	log.Println("✅ DB Connected")
	return DB
}
