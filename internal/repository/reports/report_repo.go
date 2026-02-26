package reports

import (
	"ocr-saas-backend/configs"
	"ocr-saas-backend/internal/models"

	"github.com/google/uuid"
)

func ListMyReports(
	tenantID, userID uuid.UUID,
	page, pageSize int,
) ([]models.ExpenseReport, int64, error) {

	var rows []models.ExpenseReport
	var total int64

	db := configs.DB.
		Model(&models.ExpenseReport{}).
		Where("tenant_id = ? AND user_id = ?", tenantID, userID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	err := db.
		Preload("Receipts").
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&rows).Error

	return rows, total, err
}

func GetByID(
	tenantID, id uuid.UUID,
) (*models.ExpenseReport, error) {

	var report models.ExpenseReport

	err := configs.DB.
		Preload("Receipts").
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&report).Error

	if err != nil {
		return nil, err
	}

	return &report, nil
}

func Create(report *models.ExpenseReport) error {
	return configs.DB.Create(report).Error
}

func Update(report *models.ExpenseReport) error {
	return configs.DB.Save(report).Error
}

func UpdateReportStatus(
	tenantID, reportID uuid.UUID,
	status string,
) error {
	return configs.DB.
		Model(&models.ExpenseReport{}).
		Where("id = ? AND tenant_id = ?", reportID, tenantID).
		Update("status", status).
		Error
}

func UpdateStatus(reportID uuid.UUID, status string, totalAmount int64) error {
	return configs.DB.Model(&models.ExpenseReport{}).
		Where("id = ?", reportID).
		Updates(map[string]interface{}{
			"status":       status,
			"total_amount": totalAmount,
		}).Error
}

func ListPending(
	tenantID uuid.UUID,
	page, pageSize int,
) ([]models.ExpenseReport, int64, error) {

	var rows []models.ExpenseReport
	var total int64

	db := configs.DB.
		Model(&models.ExpenseReport{}).
		Where("tenant_id = ? AND status = ?", tenantID, "SUBMITTED")

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	err := db.
		Preload("User").
		Order("created_at ASC").
		Limit(pageSize).
		Offset(offset).
		Find(&rows).Error

	return rows, total, err
}

func ListSubmitted(
	tenantID uuid.UUID,
	page, pageSize int,
) ([]models.ExpenseReport, int64, error) {

	var rows []models.ExpenseReport
	var total int64

	db := configs.DB.
		Model(&models.ExpenseReport{}).
		Where("tenant_id = ? AND status = ?", tenantID, "SUBMITTED")

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	err := db.
		Preload("User").
		Preload("Receipts").
		Order("created_at ASC").
		Limit(pageSize).
		Offset(offset).
		Find(&rows).Error

	return rows, total, err
}
