package dto

import "time"

type ExpenseReportResponse struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	TotalAmount int64             `json:"total_amount"`
	Status      string            `json:"status"`
	CreatedAt   time.Time         `json:"created_at"`
	Receipts    []ReceiptResponse `json:"receipts"`
}

type ReceiptResponse struct {
	ID          string `json:"id"`
	StoreName   string `json:"store_name"`
	TotalAmount int64  `json:"total_amount"`
	Status      string `json:"status"`
}
