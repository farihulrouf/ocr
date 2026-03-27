package dashboard

import (
	"ocr-saas-backend/internal/dto"
	"ocr-saas-backend/internal/repository/dashboard"
	"time"

	"github.com/google/uuid"
)

func GetEmployeeDashboardData(tenantID, userID uuid.UUID) (*dto.DashboardStats, error) {
	now := time.Now()
	// Default range: bulan ini
	firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// 1. Fetch Basic Stats (Receipt based)
	totalExp, _ := dashboard.GetTotalExpense(tenantID, userID, firstDayOfMonth)
	scanCount, _ := dashboard.GetScanCount(tenantID, userID)

	// PendingItems: Struk yang report_id-nya NULL (belum diajukan)
	pendingCount, _ := dashboard.GetPendingCount(tenantID, userID)

	// 2. Fetch Financial Status (Report based)
	// Menghitung nominal berdasarkan status di tabel ExpenseReport
	reportedAmt, _ := dashboard.GetAmountByReportStatus(tenantID, userID, "SUBMITTED")
	approvedAmt, _ := dashboard.GetAmountByReportStatus(tenantID, userID, "APPROVED")

	// Khusus Rejected kita ambil Count-nya (jumlah laporan yang bermasalah)
	rejectedCnt, _ := dashboard.GetRejectedCount(tenantID, userID)

	// 3. Business Logic: Saved Time (15 min per scan = 0.25 hour)
	savedTime := float64(scanCount) * 0.25

	// 4. Fetch Recent Scans (Limit 3)
	receipts, _ := dashboard.GetRecentReceipts(tenantID, userID, 3)
	var recent []dto.RecentScan
	for _, r := range receipts {
		recent = append(recent, dto.RecentScan{
			ID:        r.ID.String(),
			StoreName: r.StoreName,
			Amount:    r.TotalAmount,
			Status:    r.Status, // Mapping status OCR/Receipt
		})
	}

	// 5. Fetch History (Last 4 Months Chart)
	var history []dto.MonthlyChart
	for i := 3; i >= 0; i-- {
		t := now.AddDate(0, -i, 0)
		start := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())

		// Ambil data pengeluaran per bulan
		amount, _ := dashboard.GetTotalExpense(tenantID, userID, start)

		history = append(history, dto.MonthlyChart{
			// Format "1月", "2月" dst untuk UI Jepang
			Month:  t.Format("1月"),
			Amount: amount,
		})
	}

	// 6. Return Lengkap sesuai DTO terbaru
	return &dto.DashboardStats{
		TotalExpense:   totalExp,
		ScanCount:      scanCount,
		PendingItems:   pendingCount,
		SavedTime:      savedTime,
		ReportedAmount: reportedAmt,
		ApprovedAmount: approvedAmt,
		RejectedCount:  rejectedCnt,
		RecentScans:    recent,
		ExpenseHistory: history,
	}, nil
}
