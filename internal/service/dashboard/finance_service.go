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

	// periode bulan berjalan (dipakai tax)
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Second)

	// =========================
	// CORE FINANCE METRICS
	// =========================
	readyToPay, err := repo.GetTotalReadyToPay(tenantID)
	if err != nil {
		return nil, err
	}

	taxLiability, err := repo.GetEstTaxLiability(tenantID, startOfMonth, endOfMonth)
	if err != nil {
		return nil, err
	}

	processedTotal, err := repo.GetProcessedThisWeek(tenantID)
	if err != nil {
		return nil, err
	}

	overdueCnt, err := repo.GetOverdueCount(tenantID)
	if err != nil {
		return nil, err
	}

	// =========================
	// READY TO DISBURSE LIST
	// =========================
	payList, err := repo.GetReadyToDisburseList(tenantID, 10)
	if err != nil {
		return nil, err
	}

	// =========================
	// CASH OUTFLOW (WEEKLY)
	// =========================
	chartData, err := repo.GetWeeklyOutflow(tenantID)
	if err != nil {
		return nil, err
	}

	// fallback kalau kosong
	if len(chartData) == 0 {
		chartData = []dto.MonthlyChart{
			{Month: "Week 1", Amount: 0},
			{Month: "Week 2", Amount: 0},
			{Month: "Week 3", Amount: 0},
			{Month: "Week 4", Amount: 0},
		}
	}

	// =========================
	// RESPONSE BUILD
	// =========================
	response := &dto.FinanceDashboardStats{
		TotalReadyToPay:   readyToPay,
		OverdueCount:      overdueCnt,
		EstimatedTax:      taxLiability,
		ProcessedThisWeek: processedTotal,
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
