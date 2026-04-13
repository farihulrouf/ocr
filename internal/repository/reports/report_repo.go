package reports

import (
	"ocr-saas-backend/configs"
	"ocr-saas-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Tambahkan ini di bawah import
func selectMinimalUserInfo(db *gorm.DB) *gorm.DB {
	return db.Select("id, name, email, role")
}

func ListMyReports(
	tenantID, userID uuid.UUID, // Parameter tetap 5
	page, pageSize int,
	status string,
) ([]models.ExpenseReport, int64, error) {

	var rows []models.ExpenseReport
	var total int64

	// 1. Base Query: Selalu kunci dengan TenantID (Wajib!)
	db := configs.DB.Model(&models.ExpenseReport{}).Where("tenant_id = ?", tenantID)

	// 2. Jika userID TIDAK Nil (Berarti dia Employee), filter berdasarkan userID
	if userID != uuid.Nil {
		db = db.Where("user_id = ?", userID)
	}

	// 3. Filter Status (Optional)
	if status != "" && status != "all" {
		db = db.Where("status = ?", status)
	}

	// 4. Hitung Total & Ambil Data
	db.Count(&total)

	offset := (page - 1) * pageSize
	err := db.
		Preload("Receipts").
		Preload("User", selectMinimalUserInfo).     // EDIT DI SINI
		Preload("Approver", selectMinimalUserInfo). // EDIT DI SINI
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
		Preload("User", selectMinimalUserInfo).     // ✅ TAMBAHKAN INI (Untuk data pengaju)
		Preload("Approver", selectMinimalUserInfo). // ✅ TAMBAHKAN INI (Untuk data manajer/approver)
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
	managerID *uuid.UUID, // TAMBAHKAN PARAMETER KE-4 DI SINI
) error {
	return configs.DB.
		Model(&models.ExpenseReport{}).
		Where("id = ? AND tenant_id = ?", reportID, tenantID).
		Updates(map[string]interface{}{
			"status":         status,
			"approved_by_id": managerID, // MENGISI KOLOM PENANGGUNG JAWAB
		}).
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

func ListByStatus(
	tenantID uuid.UUID,
	page, pageSize int,
	status string,
) ([]models.ExpenseReport, int64, error) {
	var rows []models.ExpenseReport
	var total int64

	db := configs.DB.Model(&models.ExpenseReport{}).
		Where("tenant_id = ? AND status = ?", tenantID, status)

	db.Count(&total)

	offset := (page - 1) * pageSize
	err := db.
		Preload("User", selectMinimalUserInfo).     // Supaya langsing
		Preload("Approver", selectMinimalUserInfo). // Supaya langsing
		Preload("Receipts").
		Order("updated_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&rows).Error

	return rows, total, err
}
