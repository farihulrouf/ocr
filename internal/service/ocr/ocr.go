package ocr

import (
	"context"
	"fmt"
	"log"
	"ocr-saas-backend/configs"
	"ocr-saas-backend/internal/models"
	"ocr-saas-backend/internal/repository/ocr"
	"ocr-saas-backend/internal/service/ocr/aiagent"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

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
	// 🔴 TAMBAHKAN INI
	job := &models.OCRJob{
		TenantID:  tenantID,
		ReceiptID: receipt.ID,
		Status:    "PENDING",
		Engine:    "mistral",
	}

	if err := ocr.CreateOCRJob(job); err != nil {
		return nil, err
	}

	fmt.Println("[DEBUG] Receipt created:", receipt.ID)
	fmt.Println("[DEBUG] OCR Job created:", job.ID)
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

	// 1. Ambil receipt dari DB
	receipt, err := ocr.GetReceiptByID(receiptID)
	if err != nil {
		return err
	}

	// 2. Download file dari MinIO ke tmp
	//tmpPath, err := downloadReceiptFromS3(receipt.ImageURL)
	tmpPath, err := configs.S3Client.DownloadToTemp(context.Background(), receipt.ImageURL)

	if err != nil {
		receipt.Status = "FAILED"
		_ = ocr.UpdateReceipt(receipt)
		log.Println("[OCR][ERROR] Failed to download from MinIO:", err) // <-- add this
		return fmt.Errorf("failed to download from MinIO: %v", err)
	}
	log.Println("[OCR] File downloaded to:", tmpPath) // <-- add this
	//
	// hapus tmp file di akhir
	defer os.Remove(tmpPath)

	// 3. Extract OCR text
	rawText, err := ExtractText(tmpPath)
	if err != nil {
		return err
	}

	//extracted := "Full text from OCR..."
	fmt.Printf("[DEBUG] Full Text Extracted: %s\n", rawText)

	// 3.5. Ubah teks mentah menjadi JSON terstruktur (MENGGUNAKAN AI CHAT)
	// Langkah ini sangat penting agar ParseReceipt tidak error!
	structuredJSON, err := aiagent.StructureTextWithAI(rawText)
	if err != nil {
		fmt.Println("[ERROR] Structuring failed:", err)
		return err
	}

	// 4. Parse JSON hasil AI (Sekarang JSON sudah valid)
	store, total, date, taxID, isQualified, subtotal, tax, items :=
		aiagent.ParseReceipt(structuredJSON, rawText)

	fingerprint := GenerateFingerprint(store, date, total)
	existing, err := ocr.FindByFingerprint(receipt.TenantID, fingerprint)

	if err == nil && existing != nil {
		fmt.Println("[WARNING] DUPLICATE RECEIPT DETECTED:", existing.ID)

		receipt.Status = "DUPLICATE"
		receipt.OCRStatus = "COMPLETED"
		receipt.UpdatedAt = time.Now()

		_ = ocr.UpdateReceipt(receipt)

		return nil
	}
	// 5. Map ke model Receipt
	receipt.StoreName = store
	receipt.TransactionDate = &date
	receipt.TotalAmount = total
	receipt.TaxRegistrationID = taxID
	receipt.IsQualified = isQualified
	receipt.Fingerprint = fingerprint // 🔥 INI WAJIB
	receipt.OCRText = rawText         // Tetap simpan teks asli untuk audit
	receipt.OCRStatus = "COMPLETED"
	receipt.Status = "DRAFT"
	receipt.UpdatedAt = time.Now()

	// 6. Update receipt di DB
	if err := ocr.UpdateReceipt(receipt); err != nil {
		return err
	}

	// 7. Simpan item-item struk
	if len(items) > 0 {
		saveReceiptItems(receipt.ID, items, subtotal, tax)
	}

	fmt.Println("[DEBUG] OCR & Structuring completed for", receiptID)
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
	items []aiagent.ParsedItem,
	subtotal int64,
	tax int64,
) {
	fmt.Printf("[DEBUG][ITEM] Saving %d items for receipt %s\n", len(items), receiptID)

	// 1. Tentukan tax rate (8% / 10%)
	taxRate := 0
	if subtotal > 0 && tax > 0 {
		rate := float64(tax) / float64(subtotal) * 100
		if rate >= 9 {
			taxRate = 10
		} else {
			taxRate = 8
		}
	}

	var totalAllocated int64

	// 2. Loop simpan item + distribusi tax
	for i, it := range items {
		amount := it.Amount

		itemTax := int64(0)
		if subtotal > 0 && tax > 0 {
			itemTax = (amount * tax) / subtotal
			totalAllocated += itemTax
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

	// 3. Fix rounding error (biar total tax = persis)
	diff := tax - totalAllocated
	if diff != 0 && len(items) > 0 {
		fmt.Println("[DEBUG] Fix rounding tax diff:", diff)

		// ambil item terakhir
		lastItem := &models.ReceiptItem{}
		err := configs.DB.
			Where("receipt_id = ?", receiptID).
			Order("id desc").
			First(lastItem).Error

		if err == nil {
			lastItem.TaxAmount += diff
			configs.DB.Save(lastItem)

			fmt.Println("[DEBUG] Rounding diff added to last item")
		}
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

func MarkAsProcessing(receiptID string) error {
	id, _ := uuid.Parse(receiptID)
	receipt, err := ocr.GetReceiptByID(id)
	if err != nil {
		return err
	}

	receipt.Status = "PROCESSING"
	receipt.OCRStatus = "PROCESSING"
	receipt.UpdatedAt = time.Now()
	return ocr.UpdateReceipt(receipt)
}

func MarkAsSuccess(receiptID string) error {
	id, _ := uuid.Parse(receiptID)
	receipt, err := ocr.GetReceiptByID(id)
	if err != nil {
		return err
	}

	receipt.Status = "SUCCESS"
	receipt.OCRStatus = "COMPLETED"
	receipt.UpdatedAt = time.Now()
	return ocr.UpdateReceipt(receipt)
}

func MarkAsFailed(receiptID string, errMsg string) error {
	id, _ := uuid.Parse(receiptID)
	receipt, err := ocr.GetReceiptByID(id)
	if err != nil {
		return err
	}

	receipt.Status = "FAILED"
	receipt.OCRStatus = "FAILED"
	receipt.UpdatedAt = time.Now()
	return ocr.UpdateReceipt(receipt)
}

// SetOCRJobStatus - update status OCRJob dengan aman
func SetOCRJobStatus(receiptID uuid.UUID, status string, errMsg string) error {
	update := map[string]interface{}{
		"status": status,
	}

	now := time.Now()
	switch status {
	case "PROCESSING":
		update["started_at"] = now
	case "DONE", "FAILED":
		update["finished_at"] = now
	}

	if errMsg != "" {
		update["error_message"] = errMsg
	}

	return configs.DB.
		Model(&models.OCRJob{}).
		Where("receipt_id = ?", receiptID).
		Updates(update).Error
}
