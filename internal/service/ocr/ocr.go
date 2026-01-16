package ocr

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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
- parse fields (store, total, date)
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

	// 2️⃣ Extract text (OCR)
	text, err := ExtractText(receipt.ImageURL)
	if err != nil {
		fmt.Println("[ERROR] ExtractText failed:", err)
		receipt.Status = "FAILED"
		ocr.UpdateReceipt(receipt)
		return err
	}

	// 3️⃣ Parse fields dari OCR text
	store, total, date := ParseReceipt(text)
	receipt.StoreName = store
	receipt.TotalAmount = total
	receipt.TransactionDate = date

	// 4️⃣ Update DB
	receipt.OCRText = text
	receipt.OCRStatus = "COMPLETED"
	receipt.Status = "COMPLETED"
	receipt.UpdatedAt = time.Now()

	if err := ocr.UpdateReceipt(receipt); err != nil {
		fmt.Println("[ERROR] UpdateReceipt failed:", err)
		return err
	}

	fmt.Println("[DEBUG] OCR completed for", receiptID)
	return nil
}

/*
ParseReceipt
- parsing OCR text struk Jepang
- otomatis ekstrak store_name, total_amount, transaction_date
- TransactionDate dikembalikan sebagai *time.Time
*/
func ParseReceipt(text string) (storeName string, total int64, date *time.Time) {
	lines := strings.Split(text, "\n")

	// -----------------------------
	// 1️⃣ Ambil store_name
	// Cari baris yang mengandung kata "ショップ" atau huruf Jepang, sebelum alamat/TEL
	// -----------------------------
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Contoh sederhana: ambil baris yang mengandung "ショップ" atau huruf kanji/kana
		if strings.Contains(line, "ショップ") || regexp.MustCompile(`[\p{Hiragana}\p{Katakana}\p{Han}]`).MatchString(line) {
			storeName = line
			break
		}
	}
	fmt.Println("[DEBUG] Parsed store_name:", storeName)

	// -----------------------------
	// 2️⃣ Ambil total_amount
	// Cari pola "合計 ¥1,727" atau "合計 ¥1727"
	// -----------------------------
	reTotal := regexp.MustCompile(`合計\s*¥?([\d,]+)`)
	if m := reTotal.FindStringSubmatch(text); len(m) > 1 {
		totalStr := strings.ReplaceAll(m[1], ",", "")
		total, _ = strconv.ParseInt(totalStr, 10, 64)
	}
	fmt.Println("[DEBUG] Parsed total_amount:", total)

	// -----------------------------
	// 3️⃣ Ambil transaction_date
	// Cari pola YYYY年MM月DD (abaikan karakter tambahan seperti "目(水)")
	// -----------------------------
	reDate := regexp.MustCompile(`(\d{4})年(\d{1,2})月(\d{1,2})`)
	if m := reDate.FindStringSubmatch(text); len(m) > 3 {
		year, _ := strconv.Atoi(m[1])
		month, _ := strconv.Atoi(m[2])
		day, _ := strconv.Atoi(m[3])
		d := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
		date = &d
	}
	fmt.Println("[DEBUG] Parsed transaction_date:", date)

	return
}

/*
ExtractText
- dummy OCR via docker tesseract
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
