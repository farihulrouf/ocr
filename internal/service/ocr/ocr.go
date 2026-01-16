package ocr

import (
	"ocr-saas-backend/internal/models"
	"ocr-saas-backend/internal/repository/ocr"

	"github.com/google/uuid"
)

/*
UploadReceipt
- hanya perintah repo
*/
func UploadReceipt(
	tenantID uuid.UUID,
	userID uuid.UUID,
	imageURL string,
) (*models.Receipt, error) {

	receipt := &models.Receipt{
		TenantID: tenantID,
		UserID:   userID,
		ImageURL: imageURL,
		Status:   "PROCESSING",
	}

	if err := ocr.CreateReceipt(receipt); err != nil {
		return nil, err
	}

	return receipt, nil
}

/*
ProcessOCR
- async
- sekarang dummy
*/
func ProcessOCR(receiptID uuid.UUID) error {
	// nanti:
	// 1. ambil receipt
	// 2. extract text
	// 3. update DB
	return nil
}
