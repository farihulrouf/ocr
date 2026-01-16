package ocr

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"ocr-saas-backend/configs"
	"ocr-saas-backend/internal/models"
	"ocr-saas-backend/internal/repository/ocr"

	"github.com/google/uuid"
)

/*
UploadReceipt
- hanya perintah repo
*/
func UploadReceipt(tenantID, userID uuid.UUID, imageURL string) (*models.Receipt, error) {
	receipt := &models.Receipt{
		TenantID:  tenantID,
		UserID:    userID,
		ImageURL:  imageURL,
		Status:    "PROCESSING",
		OCRStatus: "PROCESSING",
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
- extract text
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

	// pastikan file ada
	if _, err := os.Stat(receipt.ImageURL); os.IsNotExist(err) {
		fmt.Println("[ERROR] File not found:", receipt.ImageURL)
		receipt.Status = "FAILED"
		ocr.UpdateReceipt(receipt)
		return fmt.Errorf("file not found: %s", receipt.ImageURL)
	}

	// 2️⃣ Extract text (dummy)
	text, err := ExtractText(receipt.ImageURL)
	if err != nil {
		fmt.Println("[ERROR] ExtractText failed:", err)
		receipt.Status = "FAILED"
		ocr.UpdateReceipt(receipt)
		return err
	}

	// 3️⃣ Update DB
	receipt.OCRText = text
	receipt.OCRStatus = "COMPLETED"
	receipt.Status = "COMPLETED"
	receipt.UpdatedAt = time.Now()

	if err := ocr.UpdateReceipt(receipt); err != nil {
		fmt.Println("[ERROR] UpdateReceipt failed:", err)
		return err
	}

	fmt.Println("[DEBUG] OCR completed for", receiptID, "text:", text)
	return nil
}

/*
ExtractText
- dummy OCR
*/
func ExtractText(imagePath string) (string, error) {
	fmt.Println("cek cek cek ")
	log.Printf("[DEBUG] OCRFromDocker called for %s", imagePath)
	filename := filepath.Base(imagePath)

	// path di dalam container
	containerPath := "/ocr/tmp/" + filename

	fmt.Println("[DEBUG] OCRFromDocker localFilePath:", imagePath)
	fmt.Println("[DEBUG] OCRFromDocker containerPath:", containerPath)

	cmd := exec.Command(
		"docker", "exec", "tesseract-ocr",
		"tesseract", containerPath, "stdout", "-l", "jpn+eng",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("docker OCR error: %w, output: %s", err, string(out))
	}
	return string(out), nil
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
