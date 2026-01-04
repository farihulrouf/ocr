package dto

import (
	"github.com/google/uuid"
)

type MyReceiptRow struct {
	ID        uuid.UUID `json:"id"`
	RecordNo  string    `json:"record_no"`
	Date      string    `json:"date"`
	StoreName string    `json:"store_name"`
	Category  string    `json:"category"`
	Taxation  string    `json:"taxation"`
	Amount    int64     `json:"amount"`
	Status    string    `json:"status"`
}
