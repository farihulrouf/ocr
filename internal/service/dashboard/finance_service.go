package dashboard

import (
	"ocr-saas-backend/internal/dto"
	// IMPORT REPOSITORY DISINI
	repo "ocr-saas-backend/internal/repository/dashboard"
	"time"

	"github.com/google/uuid"
)

func GetAdminFinanceDashboardData(tenantID uuid.UUID) (*dto.FinanceDashboardStats, error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Second)

	// Panggil menggunakan alias 'repo' agar tidak bentrok dengan package name service ini
	readyToPay, _ := repo.GetTotalReadyToPay(tenantID)
	taxLiability, _ := repo.GetEstTaxLiability(tenantID, startOfMonth, endOfMonth)
	processedWeek, _ := repo.GetProcessedThisWeek(tenantID)
	overdueCnt, _ := repo.GetOverdueCount(tenantID)

	payList, err := repo.GetReadyToDisburseList(tenantID, 10)
	if err != nil {
		return nil, err
	}

	chartData, _ := repo.GetWeeklyOutflow(tenantID)

	// Fallback Chart
	if len(chartData) == 0 {
		chartData = []dto.MonthlyChart{
			{Month: "Week 1", Amount: 0},
			{Month: "Week 2", Amount: 0},
			{Month: "Week 3", Amount: 0},
			{Month: "Week 4", Amount: 0},
		}
	}

	response := &dto.FinanceDashboardStats{
		TotalReadyToPay:   readyToPay,
		OverdueCount:      overdueCnt,
		EstimatedTax:      taxLiability,
		ProcessedThisWeek: processedWeek,
		CashOutflow:       chartData,
		ReadyToDisburse:   payList,
		SyncStatus: dto.ERPSyncStatus{
			IsConnected:   true,
			Provider:      "Xero",
			LastSync:      now.Format("02 Jan 2006 15:04"),
			PendingExport: len(payList),
		},
	}

	return response, nil
}
