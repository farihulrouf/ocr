package disbursement

import (
	"ocr-saas-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Create(tx *gorm.DB, d *models.Disbursement) error {
	return tx.Create(d).Error
}

func UpdateReportStatus(tx *gorm.DB, id uuid.UUID, status string) error {
	return tx.Model(&models.ExpenseReport{}).Where("id = ?", id).Update("status", status).Error
}

func UpdateReceiptsStatus(tx *gorm.DB, reportID uuid.UUID, status string) error {
	return tx.Model(&models.Receipt{}).Where("report_id = ?", reportID).Update("status", status).Error
}
