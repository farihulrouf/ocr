package configs

import (
	"fmt"
	"log"
	"os"

	"ocr-saas-backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB adalah variabel global agar bisa diakses oleh repository
var DB *gorm.DB

// ConnectDB melakukan inisialisasi koneksi ke database dan auto-migrate
func ConnectDB() {
	// Memastikan file .env sudah di-load (opsional jika sudah dipanggil di main)
	// godotenv.Load()

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	// Jika sslmode tidak ada di env, default ke disable
	sslmode := os.Getenv("DB_SSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}

	// Susun DSN (Database Source Name)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		host, user, password, dbname, port, sslmode,
	)

	var err error
	// Membuka koneksi ke Postgres
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})

	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	log.Println("Berhasil terhubung ke database!")

	// Jalankan AutoMigrate untuk membuat tabel berdasarkan struct models
	err = DB.AutoMigrate(
		&models.SubscriptionPlan{},
		&models.Tenant{},
		&models.CompanySetting{},
		&models.Department{},
		&models.User{},
		&models.RefreshToken{}, // <-- TAMBAHKAN INI
		&models.UserApprover{},
		&models.AccountCategory{},
		&models.TaxRate{},
		&models.PaymentMethod{},
		&models.VendorMaster{},
		&models.Receipt{},
		&models.ReceiptItem{},
		&models.ExpenseReport{},
		&models.ApprovalWorkflow{},
		&models.ApprovalStep{},
		&models.ApprovalLog{},
		&models.AuditTrail{},
		&models.ExportLog{},
	)

	if err != nil {
		log.Fatalf("Gagal migrasi database: %v", err)
	}

	log.Println("Migrasi database berhasil!")
}
