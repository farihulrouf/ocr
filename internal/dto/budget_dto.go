package dto

import (
	"github.com/google/uuid"
)

type BudgetListItem struct {
	ID              uuid.UUID `json:"id"`
	Category        string    `json:"category"`
	LimitAmount     int64     `json:"limit_amount"`
	SpentAmount     int64     `json:"spent_amount"`
	RemainingAmount int64     `json:"remaining_amount"`
	Percentage      float64   `json:"percentage"`
	IsCritical      bool      `json:"is_critical"` // True jika pemakaian > 90%
	Month           int       `json:"month"`
	Year            int       `json:"year"`
	DepartmentName  string    `json:"department_name"`
	CreatedBy       string    `json:"created_by"`
	CreatedAt       string    `json:"created_at"`
}

type FinanceBudgetSummary struct {
	TotalAllocated  int64   `json:"total_allocated"`
	TotalSpent      int64   `json:"total_spent"`
	TotalRemaining  int64   `json:"total_remaining"`
	TotalPercentage float64 `json:"total_percentage"`
	CriticalDepts   int     `json:"critical_departments_count"`
	ForecastOver    int     `json:"forecast_over_count"` // <--- TAMBAHKAN INI
	Month           int     `json:"month"`
	Year            int     `json:"year"`
}
