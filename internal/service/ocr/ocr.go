package ocr

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
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

	// 4️⃣ Parse OCR text (🔥 DARI AI)
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
	receipt.Status = "DRAFT"
	receipt.UpdatedAt = time.Now()

	// 6️⃣ Update receipt di DB
	if err := ocr.UpdateReceipt(receipt); err != nil {
		fmt.Println("[ERROR] UpdateReceipt failed:", err)
		return err
	}

	// 7️⃣ Simpan item-item struk (🔥 SOURCE = AI)
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

func saveReceiptItems(
	receiptID uuid.UUID,
	items []ParsedItem,
	subtotal int64,
	tax int64,
) {
	fmt.Printf("[DEBUG][ITEM] Saving %d items for receipt %s\n", len(items), receiptID)

	taxRate := 0
	if subtotal > 0 && tax > 0 {
		taxRate = int(float64(tax) / float64(subtotal) * 100)
	}

	for i, it := range items {
		amount := it.Amount

		itemTax := int64(0)
		if taxRate > 0 {
			itemTax = amount - (amount * 100 / int64(100+taxRate))
		}

		item := &models.ReceiptItem{
			ReceiptID:   receiptID,
			Description: it.Description,
			Amount:      amount,
			TaxAmount:   itemTax,
			TaxRate:     taxRate,
		}

		if err := ocr.CreateReceiptItem(item); err != nil {
			fmt.Printf("[ERROR][ITEM] failed save item %d: %v\n", i+1, err)
			continue
		}

		fmt.Printf(
			"[DEBUG][ITEM] saved #%d | %s | ¥%d | tax ¥%d (%d%%)\n",
			i+1, it.Description, amount, itemTax, taxRate,
		)
	}
}

func parseItemLine(line string) (description string, amount int64, ok bool) {
	re := regexp.MustCompile(`(.+?)\s*¥\s*([\d,]+)`)
	m := re.FindStringSubmatch(line)
	if len(m) < 3 {
		return "", 0, false
	}

	amountStr := strings.ReplaceAll(m[2], ",", "")
	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		return "", 0, false
	}

	description = strings.TrimSpace(m[1])
	return description, amount, true
}
