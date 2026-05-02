package ocr

import (
	"ocr-saas-backend/configs"
	"ocr-saas-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateReceipt(r *models.Receipt) error {
	return configs.DB.Create(r).Error
}

func GetReceiptByID(id uuid.UUID) (*models.Receipt, error) {
	var r models.Receipt
	if err := configs.DB.First(&r, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

func UpdateReceipt(r *models.Receipt) error {
	return configs.DB.Save(r).Error
}

// --- Fungsi untuk ReceiptItem ---

func CreateReceiptItem(item *models.ReceiptItem) error {
	return configs.DB.Create(item).Error
}

func DeleteReceiptItemsByReceiptID(receiptID uuid.UUID) error {
	return configs.DB.Where("receipt_id = ?", receiptID).Delete(&models.ReceiptItem{}).Error
}

func FindByFingerprint(
	tenantID uuid.UUID,
	fingerprint string,
) (*models.Receipt, error) {

	var receipt models.Receipt

	err := configs.DB.
		Where("tenant_id = ? AND fingerprint = ?", tenantID, fingerprint).
		First(&receipt).Error

	// 🔥 ini penting
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // ✅ tidak ada duplicate
		}
		return nil, err // ❌ error beneran
	}

	return &receipt, nil
}

func HardDeleteReceiptByID(id uuid.UUID) error {
	return configs.DB.
		Unscoped().
		Where("id = ?", id).
		Delete(&models.Receipt{}).Error
}
