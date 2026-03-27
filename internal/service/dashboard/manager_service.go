package dashboard

import (
	"ocr-saas-backend/internal/dto"
	"ocr-saas-backend/internal/repository/dashboard"
	"time"

	"github.com/google/uuid"
)

func GetManagerDashboardData(tenantID uuid.UUID) (*dto.ManagerDashboardStats, error) {
	now := time.Now()
	thisMonthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	lastMonthStart := thisMonthStart.AddDate(0, -1, 0)
	lastMonthEnd := thisMonthStart.Add(-time.Second)

	// 1. Hitung Spend & Trend
	currentSpend, _ := dashboard.GetTotalDeptExpense(tenantID, thisMonthStart, now)
	lastMonthSpend, _ := dashboard.GetTotalDeptExpense(tenantID, lastMonthStart, lastMonthEnd)

	var trend float64
	if lastMonthSpend > 0 {
		trend = (float64(currentSpend-lastMonthSpend) / float64(lastMonthSpend)) * 100
	}

	// 2. Counts
	pendingCnt, _ := dashboard.GetPendingReportsCount(tenantID)
	urgentList, _ := dashboard.GetUrgentReports(tenantID, 5)

	// 3. Chart Data (6 Months)
	var history []dto.MonthlyChart
	for i := 5; i >= 0; i-- {
		t := now.AddDate(0, -i, 0)
		start := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
		end := start.AddDate(0, 1, 0).Add(-time.Second)

		amt, _ := dashboard.GetTotalDeptExpense(tenantID, start, end)
		history = append(history, dto.MonthlyChart{
			Month:  t.Format("1月"), // Format Jepang
			Amount: amt,
		})
	}

	return &dto.ManagerDashboardStats{
		TotalDepartmentExpense: currentSpend,
		SpendTrend:             trend,
		PendingApprovalCount:   pendingCnt,
		UrgentReports:          urgentList,
		MonthlyTrend:           history,
		// Mocking Policy Alert (Bisa dihitung dari receipts yang IsQualified = false)
		PolicyAlertCount: 2,
	}, nil
}
