package dashboard

import (
	"ocr-saas-backend/configs"
	"ocr-saas-backend/internal/dto"
	"time"

	"github.com/google/uuid"
)

// GetTotalReadyToPay: Menghitung nominal dari receipt yang laporannya sudah APPROVED
func GetTotalReadyToPay(tenantID uuid.UUID) (int64, error) {
	var total int64
	// Sesuai model Mas: Nama tabel adalah expense_reports
	err := configs.DB.Table("receipts").
		Joins("JOIN expense_reports ON expense_reports.id = receipts.report_id").
		Where("receipts.tenant_id = ? AND expense_reports.status = ?", tenantID, "APPROVED").
		Select("COALESCE(SUM(receipts.total_amount), 0)").
		Scan(&total).Error
	return total, err
}

// GetOverdueCount: Menghitung invoice APPROVED yang sudah berumur > 30 hari
func GetOverdueCount(tenantID uuid.UUID) (int64, error) {
	var count int64
	deadline := time.Now().AddDate(0, 0, -30)
	err := configs.DB.Table("receipts").
		Joins("JOIN expense_reports ON expense_reports.id = receipts.report_id").
		Where("receipts.tenant_id = ? AND expense_reports.status = ? AND receipts.transaction_date < ?",
			tenantID, "APPROVED", deadline).
		Count(&count).Error
	return count, err
}

// GetEstTaxLiability: Menghitung estimasi PPN dari item di receipt yang APPROVED
func GetEstTaxLiability(tenantID uuid.UUID, start, end time.Time) (int64, error) {
	var total int64
	err := configs.DB.Table("receipt_items").
		Joins("JOIN receipts ON receipts.id = receipt_items.receipt_id").
		Joins("JOIN expense_reports ON expense_reports.id = receipts.report_id").
		Where("receipts.tenant_id = ? AND expense_reports.status = ? AND receipts.created_at BETWEEN ? AND ?",
			tenantID, "APPROVED", start, end).
		Select("COALESCE(SUM(receipt_items.tax_amount), 0)").
		Scan(&total).Error
	return total, err
}

// GetProcessedThisWeek: Mengambil data dari tabel disbursements (pencairan dana)
func GetProcessedThisWeek(tenantID uuid.UUID) (int64, error) {
	var total int64
	now := time.Now()
	// Hitung Senin minggu ini
	offset := int(now.Weekday()) - 1
	if offset < 0 {
		offset = 6
	}
	startOfWeek := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -offset)

	err := configs.DB.Table("disbursements").
		Where("tenant_id = ? AND paid_at >= ?", tenantID, startOfWeek).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	return total, err
}

// GetReadyToDisburseList: Data untuk tabel utama dashboard
func GetReadyToDisburseList(tenantID uuid.UUID, limit int) ([]dto.ReadyToPayDetail, error) {
	var results []dto.ReadyToPayDetail

	// Raw SQL dengan JOIN ke expense_reports
	err := configs.DB.Raw(`
		SELECT 
			r.id, 
			COALESCE(rep.title, r.store_name) as vendor_name, 
			COALESCE(ac.name, 'General') as category, 
			r.total_amount as amount,
			r.id::text as invoice_num,
			CASE 
				WHEN r.transaction_date < NOW() - INTERVAL '30 days' THEN 'Overdue'
				ELSE TO_CHAR(EXTRACT(DAY FROM (r.transaction_date + INTERVAL '30 days' - NOW())), '99') || ' Days Left'
			END as due_date,
			CASE 
				WHEN r.transaction_date < NOW() - INTERVAL '30 days' THEN 'Delayed'
				WHEN r.transaction_date > NOW() - INTERVAL '5 days' THEN 'Verified'
				ELSE 'Urgent'
			END as status
		FROM receipts r
		JOIN expense_reports rep ON r.report_id = rep.id
		LEFT JOIN account_categories ac ON r.account_category_id = ac.id
		WHERE r.tenant_id = ? AND rep.status = 'APPROVED'
		ORDER BY r.transaction_date ASC 
		LIMIT ?
	`, tenantID, limit).Scan(&results).Error

	return results, err
}

// GetWeeklyOutflow: Data pengeluaran bulanan (berdasarkan receipt yang statusnya PAID)
func GetWeeklyOutflow(tenantID uuid.UUID) ([]dto.MonthlyChart, error) {
	var stats []dto.MonthlyChart

	err := configs.DB.Raw(`
		SELECT 
			'Week ' || TO_CHAR(transaction_date, 'W') as month,
			SUM(total_amount) as amount
		FROM receipts
		WHERE tenant_id = ? AND status = 'PAID' AND transaction_date > NOW() - INTERVAL '3 months'
		GROUP BY month
		ORDER BY month ASC
	`, tenantID).Scan(&stats).Error

	return stats, err
}
