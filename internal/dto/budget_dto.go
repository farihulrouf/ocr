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
