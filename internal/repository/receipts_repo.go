package repository

import (
	"context"
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
			AND status = 'PROCESSING'
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

func GetAccountCategoryByID(
	tenantID uuid.UUID,
	categoryID uuid.UUID,
) (*models.AccountCategory, error) {

	var cat models.AccountCategory

	err := configs.DB.
		Where(`
			id = ?
			AND tenant_id = ?
			AND deleted_at IS NULL
		`, categoryID, tenantID).
		First(&cat).Error

	if err != nil {
		return nil, err
	}

	return &cat, nil
}

// ambil old receipt category (untuk audit)
func GetReceiptsCategorySnapshot(
	tenantID uuid.UUID,
	ids []uuid.UUID,
) ([]map[string]interface{}, error) {

	var rows []map[string]interface{}

	err := configs.DB.
		Model(&models.Receipt{}).
		Select("id, account_category_id").
		Where(`
			tenant_id = ?
			AND id IN ?
			AND deleted_at IS NULL
		`, tenantID, ids).
		Find(&rows).Error

	return rows, err
}

// bulk update
func BulkUpdateReceiptCategory(
	tenantID uuid.UUID,
	receiptIDs []uuid.UUID,
	categoryID uuid.UUID,
) (int64, error) {

	result := configs.DB.
		Model(&models.Receipt{}).
		Where(`
			tenant_id = ?
			AND id IN ?
			AND deleted_at IS NULL
		`, tenantID, receiptIDs).
		Update("account_category_id", categoryID)

	if result.Error != nil {
		return 0, result.Error
	}

	if result.RowsAffected == 0 {
		return 0, gorm.ErrRecordNotFound
	}

	return result.RowsAffected, nil
}

type ReceiptItemRepository interface {
	FindByID(ctx context.Context, id uint) (*models.ReceiptItem, error)
	Update(ctx context.Context, item *models.ReceiptItem) error
	Delete(ctx context.Context, itemID uint) error // ✅ TAMBAH INI
}

func CreateReceiptItem(
	ctx context.Context,
	item *models.ReceiptItem,
) error {

	return configs.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// 1️⃣ insert item
		if err := tx.Create(item).Error; err != nil {
			return err
		}

		// 2️⃣ hitung ulang total receipt
		var total int64
		if err := tx.
			Model(&models.ReceiptItem{}).
			Where("receipt_id = ?", item.ReceiptID).
			Select("COALESCE(SUM(amount),0)").
			Scan(&total).Error; err != nil {
			return err
		}

		// 3️⃣ update receipt.total_amount
		return tx.
			Model(&models.Receipt{}).
			Where("id = ?", item.ReceiptID).
			Update("total_amount", total).Error
	})
}

type receiptItemRepo struct{}

func NewReceiptItemRepository() ReceiptItemRepository {
	return &receiptItemRepo{}
}

func (r *receiptItemRepo) FindByID(
	ctx context.Context,
	id uint,
) (*models.ReceiptItem, error) {

	var item models.ReceiptItem
	err := configs.DB.
		WithContext(ctx).
		Preload("Receipt").
		First(&item, id).Error

	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *receiptItemRepo) Update(
	ctx context.Context,
	item *models.ReceiptItem,
) error {

	return configs.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// 1️⃣ update item
		if err := tx.
			Model(&models.ReceiptItem{}).
			Where("id = ?", item.ID).
			Updates(map[string]interface{}{
				"amount":     item.Amount,
				"updated_at": gorm.Expr("NOW()"),
			}).Error; err != nil {
			return err
		}

		// 2️⃣ hitung ulang total receipt
		var total int64
		if err := tx.
			Model(&models.ReceiptItem{}).
			Where("receipt_id = ?", item.ReceiptID).
			Select("COALESCE(SUM(amount),0)").
			Scan(&total).Error; err != nil {
			return err
		}

		// 3️⃣ update receipt.total_amount
		return tx.
			Model(&models.Receipt{}).
			Where("id = ?", item.ReceiptID).
			Update("total_amount", total).Error
	})
}

func (r *receiptItemRepo) Delete(
	ctx context.Context,
	itemID uint,
) error {

	return configs.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		var item models.ReceiptItem

		// 1️⃣ ambil item + receipt_id
		if err := tx.First(&item, itemID).Error; err != nil {
			return err
		}

		receiptID := item.ReceiptID

		// 2️⃣ delete item
		if err := tx.Delete(&models.ReceiptItem{}, itemID).Error; err != nil {
			return err
		}

		// 3️⃣ hitung ulang total receipt
		var total int64
		if err := tx.
			Model(&models.ReceiptItem{}).
			Where("receipt_id = ?", receiptID).
			Select("COALESCE(SUM(amount),0)").
			Scan(&total).Error; err != nil {
			return err
		}

		// 4️⃣ update receipt.total_amount
		return tx.
			Model(&models.Receipt{}).
			Where("id = ?", receiptID).
			Update("total_amount", total).Error
	})
}

func UpdateReceiptByID(
	tenantID, receiptID uuid.UUID,
	storeName string,
	date *time.Time,
	total *int64,
) error {

	updates := map[string]interface{}{}

	if storeName != "" {
		updates["store_name"] = storeName
	}
	if date != nil {
		updates["transaction_date"] = *date
	}
	if total != nil {
		updates["total_amount"] = *total
	}

	return configs.DB.
		Model(&models.Receipt{}).
		Where("id = ? AND tenant_id = ? AND status = 'PROCESSING'", receiptID, tenantID).
		Updates(updates).
		Error
}
