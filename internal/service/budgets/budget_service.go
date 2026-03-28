package budgets

import (
	"errors"
	"time"

	"ocr-saas-backend/configs"
	"ocr-saas-backend/internal/models"
	"ocr-saas-backend/internal/repository/budgets"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SyncBudgetSpent menghitung ulang pengeluaran APPROVED dari tabel reports ke tabel budgets
func SyncBudgetSpent(tenantID uuid.UUID, category string) error {
	now := time.Now()
	month := int(now.Month())
	year := now.Year()

	var totalSpent int64
	err := configs.DB.Model(&models.ExpenseReport{}).
		Where("tenant_id = ? AND status = ? AND EXTRACT(MONTH FROM created_at) = ? AND EXTRACT(YEAR FROM created_at) = ?",
			tenantID, "APPROVED", month, year).
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&totalSpent).Error

	if err != nil {
		return err
	}

	return configs.DB.Model(&models.Budget{}).
		Where("tenant_id = ? AND category = ? AND month = ? AND year = ?",
			tenantID, category, month, year).
		Update("spent_amount", totalSpent).Error
}

// CalculateBudgetStats mengambil data untuk Dashboard Progress Bar
func CalculateBudgetStats(tenantID uuid.UUID, category string) (map[string]interface{}, error) {
	// Jalankan sinkronisasi otomatis setiap kali data dipanggil
	_ = SyncBudgetSpent(tenantID, category)

	budget, err := budgets.GetBudgetByCategory(configs.DB, tenantID, category)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return map[string]interface{}{
				"limit_amount": 0, "spent_amount": 0, "remaining": 0, "usage_percent": 0.0, "category": category,
			}, nil
		}
		return nil, err
	}

	remaining := budget.LimitAmount - budget.SpentAmount
	usagePercent := 0.0
	if budget.LimitAmount > 0 {
		usagePercent = (float64(budget.SpentAmount) / float64(budget.LimitAmount)) * 100
	}

	return map[string]interface{}{
		"limit_amount":  budget.LimitAmount,
		"spent_amount":  budget.SpentAmount,
		"remaining":     remaining,
		"usage_percent": usagePercent,
		"category":      budget.Category,
	}, nil
}

// SetBudgetLimit menyimpan limit baru dan langsung melakukan sync
func SetBudgetLimit(tenantID, managerID uuid.UUID, category string, limit int64, month, year int) error {
	existing, err := budgets.GetBudgetByCategory(configs.DB, tenantID, category)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		newBudget := models.Budget{
			TenantID:    tenantID,
			CreatedByID: managerID,
			Category:    category,
			LimitAmount: limit,
			SpentAmount: 0,
			Month:       month,
			Year:        year,
		}
		if err := budgets.CreateOrUpdateBudget(configs.DB, &newBudget); err != nil {
			return err
		}
		return SyncBudgetSpent(tenantID, category)
	}

	if err != nil {
		return err
	}

	existing.LimitAmount = limit
	existing.CreatedByID = managerID
	if err := budgets.CreateOrUpdateBudget(configs.DB, existing); err != nil {
		return err
	}
	return SyncBudgetSpent(tenantID, category)
}

// GetTenantBudgets mengambil history budget
func GetTenantBudgets(tenantID uuid.UUID, year int) ([]models.Budget, error) {
	var rows []models.Budget
	db := configs.DB.Where("tenant_id = ?", tenantID)
	if year != 0 {
		db = db.Where("year = ?", year)
	}
	err := db.Order("year DESC, month DESC").Find(&rows).Error
	return rows, err
}

// ConsumeBudgetLogic dipanggil di Report Service saat Approval
func ConsumeBudgetLogic(tx *gorm.DB, tenantID uuid.UUID, amount int64) error {
	budget, err := budgets.GetBudgetByCategory(tx, tenantID, "General")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	if budget.SpentAmount+amount > budget.LimitAmount {
		return errors.New("limit budget tidak mencukupi")
	}

	return budgets.UpdateSpentAmount(tx, budget.ID, amount)
}
