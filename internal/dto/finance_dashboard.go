package dto

import (
	"github.com/google/uuid"
)

// FinanceDashboardStats menggunakan MonthlyChart yang SUDAH ADA di dashboard.go
type FinanceDashboardStats struct {
	TotalReadyToPay   int64              `json:"total_ready_to_pay"`
	OverdueCount      int64              `json:"overdue_count"`
	EstimatedTax      int64              `json:"estimated_tax"`
	ProcessedThisWeek int64              `json:"processed_this_week"`
	CashOutflow       []MonthlyChart     `json:"cash_outflow"` // <--- Ini AMAN, otomatis pakai yang sudah ada
	ReadyToDisburse   []ReadyToPayDetail `json:"ready_to_disburse"`
	SyncStatus        ERPSyncStatus      `json:"sync_status"`
}

type ReadyToPayDetail struct {
	ID         uuid.UUID `json:"id"`
	VendorName string    `json:"vendor_name"`
	Category   string    `json:"category"`
	Amount     int64     `json:"amount"`
	DueDate    string    `json:"due_date"`
	Status     string    `json:"status"`
	InvoiceNum string    `json:"invoice_num"`
}

type ERPSyncStatus struct {
	IsConnected   bool   `json:"is_connected"`
	Provider      string `json:"provider"`
	LastSync      string `json:"last_sync"`
	PendingExport int    `json:"pending_export"`
}
