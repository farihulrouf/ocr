package repository

import (
	"ocr-saas-backend/configs"
	"ocr-saas-backend/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetMyReceipts(
	tenantID, userID uuid.UUID,
	page, pageSize int,
	q, status, sort string,
) ([]models.Receipt, int64, error) {

	var receipts []models.Receipt
	var total int64

	db := configs.DB.
		Model(&models.Receipt{}).
		Preload("AccountCategory").
		Where("tenant_id = ? AND user_id = ?", tenantID, userID)

	// 🔍 search store name
	if q != "" {
		db = db.Where("store_name ILIKE ?", "%"+q+"%")
	}

	// 🟡 filter status
	if status != "" {
		db = db.Where("status = ?", status)
	}

	// count
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	// sorting
	if sort != "" {
		db = db.Order(sort)
	} else {
		db = db.Order("transaction_date DESC")
	}

	err := db.
		Limit(pageSize).
		Offset(offset).
		Find(&receipts).Error

	return receipts, total, err
}

func GetAllReceipts(
	tenantID uuid.UUID,
	page, pageSize int,
	q, status, sort string,
) ([]models.Receipt, int64, error) {

	var receipts []models.Receipt
	var total int64

	db := configs.DB.
		Model(&models.Receipt{}).
		Preload("AccountCategory").
		Preload("User").
		Where("tenant_id = ?", tenantID)

	// 🔍 search store name
	if q != "" {
		db = db.Where("store_name ILIKE ?", "%"+q+"%")
	}

	// 🟡 filter status
	if status != "" {
		db = db.Where("status = ?", status)
	}

	// count
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	// sorting
	if sort != "" {
		db = db.Order(sort)
	} else {
		db = db.Order("transaction_date DESC")
	}

	err := db.
		Limit(pageSize).
		Offset(offset).
		Find(&receipts).Error

	return receipts, total, err
}

func GetReceiptDetailByID(
	tenantID uuid.UUID,
	receiptID uuid.UUID,
) (*models.Receipt, error) {

	var receipt models.Receipt

	err := configs.DB.
		Model(&models.Receipt{}).
		Preload("User").
		Preload("AccountCategory").
		Preload("LineItems").
		Where("id = ? AND tenant_id = ?", receiptID, tenantID).
		First(&receipt).Error

	if err != nil {
		return nil, err
	}

	return &receipt, nil
}

func ConfirmReceiptByID(
	tenantID uuid.UUID,
	receiptID uuid.UUID,
	total int64,
	date time.Time,
) error {

	return configs.DB.
		Model(&models.Receipt{}).
		Where(`
			id = ? 
			AND tenant_id = ?
			AND status = 'PENDING'
		`, receiptID, tenantID).
		Updates(map[string]interface{}{
			"total_amount":     total,
			"transaction_date": date,
			"status":           "APPROVED",
		}).Error
}

func DeleteReceiptByIDManager(
	tenantID uuid.UUID,
	receiptID uuid.UUID,
) error {

	result := configs.DB.
		Where(`
			id = ?
			AND tenant_id = ?
		`, receiptID, tenantID).
		Delete(&models.Receipt{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func BulkDeleteReceiptsByManager(
	tenantID uuid.UUID,
	ids []uuid.UUID,
) (int64, error) {

	result := configs.DB.
		Where(`
			id IN ?
			AND tenant_id = ?
			AND status = 'PENDING'
		`, ids, tenantID).
		Delete(&models.Receipt{})

	if result.Error != nil {
		return 0, result.Error
	}

	if result.RowsAffected == 0 {
		return 0, gorm.ErrRecordNotFound
	}

	return result.RowsAffected, nil
}

func BulkRestoreReceiptsByIDs(
	tenantID uuid.UUID,
	ids []uuid.UUID,
) (int64, error) {

	result := configs.DB.
		Unscoped().
		Model(&models.Receipt{}).
		Where(`
			id IN ?
			AND tenant_id = ?
			AND deleted_at IS NOT NULL
		`, ids, tenantID).
		Update("deleted_at", nil)

	if result.Error != nil {
		return 0, result.Error
	}

	if result.RowsAffected == 0 {
		return 0, gorm.ErrRecordNotFound
	}

	return result.RowsAffected, nil
}

// internal/repository/receipt_repository.go
func BulkUpdateReceiptStatusTx(
	tx *gorm.DB,
	tenantID uuid.UUID,
	ids []uuid.UUID,
	newStatus string,
) (int64, error) {

	result := tx.Model(&models.Receipt{}).
		Where(`
			id IN ?
			AND tenant_id = ?
			AND status = 'PENDING'
			AND deleted_at IS NULL
		`, ids, tenantID).
		Update("status", newStatus)

	if result.Error != nil {
		return 0, result.Error
	}

	if result.RowsAffected == 0 {
		return 0, gorm.ErrRecordNotFound
	}

	return result.RowsAffected, nil
}
