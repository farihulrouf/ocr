package dashboard

import (
	"ocr-saas-backend/configs"
	"ocr-saas-backend/internal/models"
	"time"

	"github.com/google/uuid"
)

// --- EXISTING FUNCTIONS ---

func GetTotalExpense(tenantID, userID uuid.UUID, since time.Time) (int64, error) {
	var total int64
	err := configs.DB.Model(&models.Receipt{}).
		Where("tenant_id = ? AND user_id = ? AND created_at >= ? AND status != ?", tenantID, userID, since, "REJECTED").
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&total).Error
	return total, err
}

func GetScanCount(tenantID, userID uuid.UUID) (int64, error) {
	var count int64
	err := configs.DB.Model(&models.Receipt{}).
		Where("tenant_id = ? AND user_id = ?", tenantID, userID).
		Count(&count).Error
	return count, err
}

// GetPendingCount: Struk yang baru di-upload tapi belum masuk ke report manapun
func GetPendingCount(tenantID, userID uuid.UUID) (int64, error) {
	var count int64
	err := configs.DB.Model(&models.Receipt{}).
		Where("tenant_id = ? AND user_id = ? AND report_id IS NULL", tenantID, userID).
		Count(&count).Error
	return count, err
}

func GetRecentReceipts(tenantID, userID uuid.UUID, limit int) ([]models.Receipt, error) {
	var receipts []models.Receipt
	err := configs.DB.Where("tenant_id = ? AND user_id = ?", tenantID, userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&receipts).Error
	return receipts, err
}

// --- NEW FUNCTIONS (REPORT STATUS) ---

// GetAmountByReportStatus: Menghitung total nominal berdasarkan status report
// Status: "SUBMITTED" (Sedang diajukan), "APPROVED" (Sudah cair), "REJECTED" (Ditolak)
func GetAmountByReportStatus(tenantID, userID uuid.UUID, status string) (int64, error) {
	var total int64
	err := configs.DB.Model(&models.ExpenseReport{}).
		Where("tenant_id = ? AND user_id = ? AND status = ?", tenantID, userID, status).
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&total).Error
	return total, err
}

// GetRejectedCount: Menghitung jumlah laporan yang ditolak (Count, bukan SUM)
func GetRejectedCount(tenantID, userID uuid.UUID) (int64, error) {
	var count int64
	err := configs.DB.Model(&models.ExpenseReport{}).
		Where("tenant_id = ? AND user_id = ? AND status = ?", tenantID, userID, "REJECTED").
		Count(&count).Error
	return count, err
}
