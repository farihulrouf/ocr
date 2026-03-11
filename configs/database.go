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

	// Ambil *sql.DB dari gorm.DB untuk atur pool
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Gagal ambil sql.DB: %v", err)
	}

	// Atur connection pool
	sqlDB.SetMaxOpenConns(20)           // Maks koneksi aktif
	sqlDB.SetMaxIdleConns(10)           // Maks koneksi idle
	sqlDB.SetConnMaxLifetime(time.Hour) // Lama koneksi sebelum diganti

	log.Println("✅ Connection pool siap digunakan")

	// Jalankan AutoMigrate untuk membuat tabel berdasarkan struct models
	err = DB.AutoMigrate(
		&models.SubscriptionPlan{},
		&models.Tenant{},
		&models.CompanySetting{},
		&models.Department{},
		&models.User{},
		&models.RefreshToken{}, // Token login
		&models.UserApprover{}, // Relasi bawahan → atasan
		&models.AccountCategory{},
		&models.TaxRate{},
		&models.PaymentMethod{},
		&models.VendorMaster{},

		// ======= IMPORTANT: Parent dulu baru Child ========
		&models.ExpenseReport{},    // Parent Receipt, ApprovalLog, ExportLog
		&models.ApprovalWorkflow{}, // Parent ApprovalStep
		&models.ApprovalStep{},     // Child

		&models.Receipt{},     // Child dari ExpenseReport
		&models.ReceiptItem{}, // Child dari Receipt

		&models.ApprovalLog{}, // Child ExpenseReport
		&models.AuditTrail{},  // Child Tenant/User
		&models.ExportLog{},   // Child ExpenseReport
		&models.TenantUsage{}, // Statistik pemakaian
		&models.OCRJob{},
	)

	if err != nil {
		log.Fatalf("Gagal migrasi database: %v", err)
	}

	log.Println("Migrasi database berhasil!")
}
