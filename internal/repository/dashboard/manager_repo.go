package dashboard

import (
	"ocr-saas-backend/configs"
	"ocr-saas-backend/internal/dto"
	"ocr-saas-backend/internal/models"
	"time"

	"github.com/google/uuid"
)

// GetTotalDeptExpense: Total spend untuk seluruh tenant/dept (Status APPROVED)
func GetTotalDeptExpense(tenantID uuid.UUID, start, end time.Time) (int64, error) {
	var total int64
	err := configs.DB.Model(&models.ExpenseReport{}).
		Where("tenant_id = ? AND status = ? AND created_at BETWEEN ? AND ?", tenantID, "APPROVED", start, end).
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&total).Error
	return total, err
}

// GetPendingReports: Laporan yang butuh approval
func GetPendingReportsCount(tenantID uuid.UUID) (int64, error) {
	var count int64
	err := configs.DB.Model(&models.ExpenseReport{}).
		Where("tenant_id = ? AND status = ?", tenantID, "SUBMITTED").
		Count(&count).Error
	return count, err
}

// GetUrgentReports: List laporan yang sudah "mengendap" lama
func GetUrgentReports(tenantID uuid.UUID, limit int) ([]dto.UrgentReport, error) {
	var results []dto.UrgentReport
	// Query manual untuk hitung selisih hari (Postgres version)
	err := configs.DB.Raw(`
        SELECT r.id, u.name as staff_name, r.title, r.total_amount as amount, r.status,
        EXTRACT(DAY FROM (NOW() - r.created_at)) as days_pending
        FROM expense_reports r
        JOIN users u ON r.user_id = u.id
        WHERE r.tenant_id = ? AND r.status = 'SUBMITTED'
        ORDER BY days_pending DESC
        LIMIT ?
    `, tenantID, limit).Scan(&results).Error
	return results, err
}
