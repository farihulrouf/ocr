package repository

import (
	"ocr-saas-backend/configs"
	"ocr-saas-backend/internal/models"

	"github.com/google/uuid"
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
