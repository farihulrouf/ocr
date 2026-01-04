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

type ReceiptUserInfo struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Name  string    `json:"name"`
}

type AdminReceiptRow struct {
	ID        uuid.UUID       `json:"id"`
	RecordNo  string          `json:"record_no"`
	Date      string          `json:"date"`
	StoreName string          `json:"store_name"`
	Category  string          `json:"category"`
	Taxation  string          `json:"taxation"`
	Amount    int64           `json:"amount"`
	Status    string          `json:"status"`
	User      ReceiptUserInfo `json:"user"`
}
