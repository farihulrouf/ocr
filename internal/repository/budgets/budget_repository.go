package budgets

import (
	"ocr-saas-backend/configs"
	"ocr-saas-backend/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GetBudgetByCategory mengambil data budget berdasarkan tenant dan kategori
func GetBudgetByCategory(db *gorm.DB, tenantID uuid.UUID, category string) (*models.Budget, error) {
	var budget models.Budget
	// Tambahkan .Preload("Department") dan .Preload("Creator")
	err := db.Preload("Department").
		Preload("Creator").
		Where("tenant_id = ? AND category = ?", tenantID, category).
		First(&budget).Error
	return &budget, err
}

// CreateOrUpdateBudget menggunakan GORM Save untuk Upsert berdasarkan ID atau manual Logic
func CreateOrUpdateBudget(db *gorm.DB, budget *models.Budget) error {
	// Jika ID kosong, GORM akan Create. Jika ada, akan Update.
	return db.Save(budget).Error
}

// UpdateSpentAmount khusus untuk menambah pemakaian budget secara atomik
func UpdateSpentAmount(tx *gorm.DB, budgetID uuid.UUID, amount int64) error {
	return tx.Model(&models.Budget{}).
		Where("id = ?", budgetID).
		UpdateColumn("spent_amount", gorm.Expr("spent_amount + ?", amount)).Error
}

// internal/service/budgets/service.go

func SyncBudgetSpent(tenantID uuid.UUID, category string) error {
	now := time.Now()
	month := int(now.Month())
	year := now.Year()

	// 1. Hitung total realita pengeluaran dari tabel ExpenseReports
	var totalSpent int64
	err := configs.DB.Model(&models.ExpenseReport{}).
		Where("tenant_id = ? AND status = ? AND EXTRACT(MONTH FROM created_at) = ? AND EXTRACT(YEAR FROM created_at) = ?",
			tenantID, "APPROVED", month, year).
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&totalSpent).Error

	if err != nil {
		return err
	}

	// 2. Update kolom spent_amount di tabel budgets
	return configs.DB.Model(&models.Budget{}).
		Where("tenant_id = ? AND category = ? AND month = ? AND year = ?",
			tenantID, category, month, year).
		Update("spent_amount", totalSpent).Error
}
