package dashboard

import (
	"ocr-saas-backend/configs"
	"ocr-saas-backend/internal/models"
	"time"

	"github.com/google/uuid"
)

// Get recent N scans for a tenant
func ListRecentScans(tenantID uuid.UUID, limit int) ([]models.Receipt, error) {
	var receipts []models.Receipt
	err := configs.DB.
		Preload("User").
		Where("tenant_id = ?", tenantID).
		Order("created_at DESC").
		Limit(limit).
		Find(&receipts).Error
	return receipts, err
}

// Get receipts in a date range
func ListReceiptsByDateRange(tenantID uuid.UUID, start, end time.Time) ([]models.Receipt, error) {
	var receipts []models.Receipt
	err := configs.DB.
		Where("tenant_id = ? AND transaction_date BETWEEN ? AND ?", tenantID, start, end).
		Order("transaction_date ASC").
		Find(&receipts).Error
	return receipts, err
}

// Optional: Group by category with sum
func GetCategoryStats(tenantID uuid.UUID) (map[string]int64, error) {
	type result struct {
		Category string
		Total    int64
	}
	var res []result

	err := configs.DB.
		Model(&models.Receipt{}).
		Select("account_category_id as category, SUM(total_amount) as total").
		Where("tenant_id = ?", tenantID).
		Group("account_category_id").
		Scan(&res).Error

	if err != nil {
		return nil, err
	}

	stats := make(map[string]int64)
	for _, r := range res {
		stats[r.Category] = r.Total
	}

	return stats, nil
}
