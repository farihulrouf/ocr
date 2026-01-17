package ocr

import (
	"fmt"
	"os"
	"time"

	"ocr-saas-backend/configs"
	"ocr-saas-backend/internal/models"
	"ocr-saas-backend/internal/repository/ocr"

	"github.com/google/uuid"
)

/*
UploadReceipt
- hanya membuat record receipt di DB
*/
func UploadReceipt(tenantID, userID uuid.UUID, imageURL string) (*models.Receipt, error) {
	receipt := &models.Receipt{
		TenantID:  tenantID,
		UserID:    userID,
		ImageURL:  imageURL,
		Status:    "PROCESSING",
		OCRStatus: "PROCESSING",
		// CreatedAt dan UpdatedAt otomatis dari Base.BeforeCreate
	}

	if err := ocr.CreateReceipt(receipt); err != nil {
		return nil, err
	}
	fmt.Println("[DEBUG] Receipt created:", receipt.ID)
	return receipt, nil
}

/*
PushToQueue
- masukkan receipt ID ke Redis queue
*/
func PushToQueue(receiptID uuid.UUID) error {
	fmt.Println("[DEBUG] Pushing to queue:", receiptID)
	return configs.RedisClient.LPush(configs.Ctx, "ocr:queue", receiptID.String()).Err()
}

/*
ProcessOCR
- ambil receipt
- extract OCR text
- parse fields (store, total, date, tax_id, etc)
- update DB
*/
func ProcessOCR(receiptID uuid.UUID) error {
	fmt.Println("[DEBUG] Starting OCR for receipt:", receiptID)

	// 1️⃣ Ambil receipt dari DB
	receipt, err := ocr.GetReceiptByID(receiptID)
	if err != nil {
		fmt.Println("[ERROR] GetReceiptByID failed:", err)
		return err
	}
	fmt.Println("[DEBUG] Receipt fetched:", receipt.ID, receipt.ImageURL)

	// 2️⃣ Pastikan file ada
	if _, err := os.Stat(receipt.ImageURL); os.IsNotExist(err) {
		fmt.Println("[ERROR] File not found:", receipt.ImageURL)

		receipt.Status = "FAILED"
		receipt.OCRStatus = "FAILED"
		receipt.UpdatedAt = time.Now()
		_ = ocr.UpdateReceipt(receipt)

		return fmt.Errorf("file not found: %s", receipt.ImageURL)
	}

	// 3️⃣ Extract OCR text
	text, err := ExtractText(receipt.ImageURL)
	if err != nil {
		fmt.Println("[ERROR] ExtractText failed:", err)

		receipt.Status = "FAILED"
		receipt.OCRStatus = "FAILED"
		receipt.UpdatedAt = time.Now()
		_ = ocr.UpdateReceipt(receipt)

		return err
	}

	// 4️⃣ Parse OCR text
	store, total, date, taxID, isQualified, subtotal, tax, items :=
		ParseReceipt(text)

	// 5️⃣ Map hasil parse ke model Receipt
	receipt.StoreName = store
	receipt.TransactionDate = date
	receipt.TotalAmount = total
	receipt.TaxRegistrationID = taxID
	receipt.IsQualified = isQualified

	receipt.OCRText = text
	receipt.OCRStatus = "COMPLETED"
	receipt.Status = "COMPLETED"
	receipt.UpdatedAt = time.Now()

	// 6️⃣ Update receipt di DB
	if err := ocr.UpdateReceipt(receipt); err != nil {
		fmt.Println("[ERROR] UpdateReceipt failed:", err)
		return err
	}

	// 7️⃣ Simpan item-item struk
	if len(items) > 0 {
		saveReceiptItems(receipt.ID, items, subtotal, tax)
	}

	fmt.Println("[DEBUG] OCR completed for", receiptID)
	return nil
}

/*
ProcessOCRString
- helper untuk worker (string -> uuid)
*/
func ProcessOCRString(receiptID string) error {
	id, err := uuid.Parse(receiptID)
	if err != nil {
		fmt.Println("[ERROR] Invalid UUID:", receiptID)
		return err
	}
	return ProcessOCR(id)
}

/*
saveReceiptItems - simpan items ke tabel receipt_items
*/
func saveReceiptItems(receiptID uuid.UUID, items []string, subtotal int64, tax int64) {
	fmt.Printf("[DEBUG] Saving %d items for receipt %s\n", len(items), receiptID)

	// Hapus items lama jika ada
	/*if err := ocr.DeleteReceiptItemsByReceiptID(receiptID); err != nil {
		fmt.Printf("[WARNING] Failed to delete old items: %v\n", err)
	}

	// Hitung rata-rata tax rate jika ada subtotal dan tax
	taxRate := 0
	if subtotal > 0 && tax > 0 {
		taxRate = int(float64(tax) / float64(subtotal) * 100)
		fmt.Printf("[DEBUG] Calculated tax rate: %d%%\n", taxRate)
	}

	// Parse dan simpan setiap item
	var receiptItems []models.ReceiptItem
	for i, itemStr := range items {
		description, amount := parseItemLine(itemStr)

		// Hitung tax amount per item
		taxAmount := int64(0)
		if taxRate > 0 && amount > 0 {
			// Asumsi tax termasuk dalam amount (Japanese style)
			taxAmount = amount - (amount * 100 / int64(100+taxRate))
		}

		receiptItem := models.ReceiptItem{
			ReceiptID:   receiptID,
			Description: description,
			Amount:      amount,
			TaxAmount:   taxAmount,
			TaxRate:     taxRate,
		}

		receiptItems = append(receiptItems, receiptItem)
		fmt.Printf("  Item %d: %s - ¥%d (Tax: ¥%d, Rate: %d%%)\n",
			i+1, description, amount, taxAmount, taxRate)
	}

	// Simpan ke database
	if len(receiptItems) > 0 {
		if err := ocr.CreateReceiptItems(receiptItems); err != nil {
			fmt.Printf("[ERROR] Failed to save receipt items: %v\n", err)
		} else {
			fmt.Printf("[DEBUG] Successfully saved %d receipt items\n", len(receiptItems))
		}
	}
	*/
}
