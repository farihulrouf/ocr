package dto

type DashboardStats struct {
	TotalExpense int64   `json:"total_expense"`
	ScanCount    int64   `json:"scan_count"`
	PendingItems int64   `json:"pending_items"`
	SavedTime    float64 `json:"saved_time"`

	// --- STATUS REAL-TIME KEUANGAN ---
	ReportedAmount int64 `json:"reported_amount"` // Sedang diajukan (Wait Manager)
	ApprovedAmount int64 `json:"approved_amount"` // Sudah disetujui (Duit aman)
	//RejectedAmount int64 `json:"rejected_amount"` // Ditolak (Butuh revisi)
	RejectedCount int64 `json:"rejected_count"`

	ExpenseHistory []MonthlyChart `json:"expense_history"`
	RecentScans    []RecentScan   `json:"recent_scans"`
}

type MonthlyChart struct {
	Month  string `json:"month"`
	Amount int64  `json:"amount"`
}

type RecentScan struct {
	ID        string `json:"id"`
	StoreName string `json:"store_name"`
	Amount    int64  `json:"amount"`
	Status    string `json:"status"`
}
