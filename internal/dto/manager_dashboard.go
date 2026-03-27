package dto

type ManagerDashboardStats struct {
	// --- TOP ROW STATS ---
	TotalDepartmentExpense int64   `json:"total_dept_expense"`
	SpendTrend             float64 `json:"spend_trend"` // Persentase vs bulan lalu
	PendingApprovalCount   int64   `json:"pending_approvals"`
	OverdueCount           int64   `json:"overdue_count"` // Pending > 3 hari
	PolicyAlertCount       int64   `json:"policy_alerts"` // Misal: Struk tanpa Tax ID

	// --- CHART DATA ---
	MonthlyTrend []MonthlyChart `json:"monthly_trend"` // Perbandingan 4-6 bulan terakhir

	// --- URGENT ACTIONS (LIST) ---
	UrgentReports []UrgentReport `json:"urgent_reports"`
}

type UrgentReport struct {
	ID          string `json:"id"`
	StaffName   string `json:"staff_name"`
	Title       string `json:"title"`
	Amount      int64  `json:"amount"`
	DaysPending int    `json:"days_pending"`
	Status      string `json:"status"`
}
