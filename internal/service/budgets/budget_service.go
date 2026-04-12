package budgets

import (
	"errors"
	"fmt"
	"math"
	"time"

	"ocr-saas-backend/configs"
	"ocr-saas-backend/internal/dto"
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
	// COALESCE digunakan agar jika tidak ada data, return 0 bukan NULL
	err := configs.DB.Model(&models.ExpenseReport{}).
		Where("tenant_id = ? AND status = ? AND EXTRACT(MONTH FROM created_at) = ? AND EXTRACT(YEAR FROM created_at) = ?",
			tenantID, "APPROVED", month, year).
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&totalSpent).Error

	if err != nil {
		return err
	}

	// Update ke table budgets
	return configs.DB.Model(&models.Budget{}).
		Where("tenant_id = ? AND category = ? AND month = ? AND year = ?",
			tenantID, category, month, year).
		Update("spent_amount", totalSpent).Error
}

// CalculateBudgetStats adalah fungsi yang dipanggil oleh API Controller
// untuk mengisi data di Dashboard Manager.
func CalculateBudgetStats(tenantID uuid.UUID, category string) (map[string]interface{}, error) {
	// 1. Sinkronisasi data agar angka 'spent_amount' terbaru dari database
	_ = SyncBudgetSpent(tenantID, category)

	// 2. Ambil data budget dari DB
	budget, err := budgets.GetBudgetByCategory(configs.DB, tenantID, category)
	if err != nil {
		// Handle error jika budget belum di-set
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return map[string]interface{}{"limit_amount": 0, "spent_amount": 0}, nil
		}
		return nil, err
	}

	// --- HITUNG BASIC ---
	remaining := budget.LimitAmount - budget.SpentAmount
	usagePercent := 0.0
	if budget.LimitAmount > 0 {
		usagePercent = (float64(budget.SpentAmount) / float64(budget.LimitAmount)) * 100
	}

	// --- TENTUKAN WARNA (SAFE/WARNING/CRITICAL) ---
	statusLevel := "safe"
	if usagePercent >= 90 {
		statusLevel = "critical"
	} else if usagePercent >= 70 {
		statusLevel = "warning"
	}

	// ========================================================
	// DISINI TEMPAT KODE FORECASTING YANG MAS TANYAKAN TADI:
	// ========================================================
	now := time.Now()
	// 1. Ambil jumlah hari dalam bulan ini
	daysInMonth := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, time.Local).Day()
	currentDay := now.Day()

	// 2. Gunakan float64 eksplisit agar tidak hasil 0
	spentFloat := float64(budget.SpentAmount)
	daysInMonthFloat := float64(daysInMonth)
	currentDayFloat := float64(currentDay)

	estimatedEndMonth := 0.0
	avgPerDay := 0.0

	if currentDayFloat > 0 {
		// Logikanya: Jika 28 hari sudah pakai 99k, maka rata-rata per hari = 99k / 28.
		// Lalu diproyeksikan ke total hari (misal 31 hari).
		avgPerDay = spentFloat / currentDayFloat
		estimatedEndMonth = math.Round(avgPerDay * daysInMonthFloat)
	}

	// 3. Cek apakah prediksi akhir bulan bakal jebol
	isOverForecast := false
	if estimatedEndMonth > float64(budget.LimitAmount) && budget.LimitAmount > 0 {
		isOverForecast = true
	}

	// --- DEBUG TERMINAL ---
	fmt.Printf("\n[DEBUG SEIDO] Spent: %.2f | Day: %.0f/%v | Forecast: %.2f\n",
		spentFloat, currentDayFloat, daysInMonth, estimatedEndMonth)

	// 4. Kirim hasilnya ke Frontend (React)
	return map[string]interface{}{
		"limit_amount":        budget.LimitAmount,
		"spent_amount":        budget.SpentAmount,
		"remaining":           remaining,
		"usage_percent":       math.Round(usagePercent*10) / 10,
		"category":            budget.Category,
		"status_level":        statusLevel,
		"estimated_end_month": estimatedEndMonth, // <--- Ini yang bikin angka ¥0 jadi ¥100rb+
		"is_over_forecast":    isOverForecast,
		"current_day":         currentDay,
		"days_in_month":       daysInMonth,
	}, nil
	// ========================================================
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
func GetTenantBudgets(tID uuid.UUID, year int) ([]dto.BudgetListItem, error) {
	var rows []models.Budget

	// 1. Ambil data dari DB dengan Preload
	err := configs.DB.
		Preload("Department").
		Preload("Creator").
		Where("tenant_id = ?", tID).
		Where("year = ?", year).
		Order("month DESC, created_at DESC").
		Find(&rows).Error

	if err != nil {
		return nil, err
	}

	// 2. Mapping ke DTO
	var response []dto.BudgetListItem
	for _, b := range rows {
		// Logic nama departemen
		deptName := "General"
		if b.Department != nil {
			deptName = b.Department.Name
		}

		// Logic Persentase (Rounding ke 2 desimal)
		var percent float64
		if b.LimitAmount > 0 {
			rawPercent := (float64(b.SpentAmount) / float64(b.LimitAmount)) * 100
			percent = math.Round(rawPercent*100) / 100
		}

		// Logic Alert Level (Kritis jika > 90%)
		isCritical := percent >= 90.0

		// Masukkan ke array response
		response = append(response, dto.BudgetListItem{
			ID:              b.ID,
			Category:        b.Category,
			LimitAmount:     b.LimitAmount,
			SpentAmount:     b.SpentAmount,
			RemainingAmount: b.LimitAmount - b.SpentAmount,
			Percentage:      percent,
			IsCritical:      isCritical,
			Month:           b.Month,
			Year:            b.Year,
			DepartmentName:  deptName,
			CreatedBy:       b.Creator.Name,
			CreatedAt:       b.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return response, nil
}

// ConsumeBudgetLogic dipanggil di Report Service saat Approval (Manual Update)
func ConsumeBudgetLogic(tx *gorm.DB, tenantID uuid.UUID, amount int64) error {
	// Default kategori 'General'
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
