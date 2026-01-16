package ocr

import (
	"fmt"
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

	// pastikan file ada
	if _, err := os.Stat(receipt.ImageURL); os.IsNotExist(err) {
		fmt.Println("[ERROR] File not found:", receipt.ImageURL)
		receipt.Status = "FAILED"
		receipt.OCRStatus = "FAILED"
		receipt.UpdatedAt = time.Now()
		ocr.UpdateReceipt(receipt)
		return fmt.Errorf("file not found: %s", receipt.ImageURL)
	}

	// 2️⃣ Extract text (OCR)
	text, err := ExtractText(receipt.ImageURL)
	if err != nil {
		fmt.Println("[ERROR] ExtractText failed:", err)
		receipt.Status = "FAILED"
		receipt.OCRStatus = "FAILED"
		receipt.UpdatedAt = time.Now()
		ocr.UpdateReceipt(receipt)
		return err
	}

	// 3️⃣ Parse semua fields dari OCR text
	store, total, date, taxID, isQualified, subtotal, tax, items := ParseReceipt(text)

	// Set semua field ke receipt
	receipt.StoreName = store
	receipt.TotalAmount = total
	receipt.TransactionDate = date
	receipt.TaxRegistrationID = taxID
	receipt.IsQualified = isQualified

	// Simpan OCR text dan status
	receipt.OCRText = text
	receipt.OCRStatus = "COMPLETED"
	receipt.Status = "COMPLETED"
	receipt.UpdatedAt = time.Now()

	// 4️⃣ Update DB receipt
	if err := ocr.UpdateReceipt(receipt); err != nil {
		fmt.Println("[ERROR] UpdateReceipt failed:", err)
		return err
	}

	// 5️⃣ Simpan items ke tabel receipt_items
	if len(items) > 0 {
		saveReceiptItems(receipt.ID, items, subtotal, tax)
	}

	fmt.Println("[DEBUG] OCR completed for", receiptID)
	return nil
}

/*
ParseReceipt
- parsing OCR text struk Jepang
- otomatis ekstrak semua field yang ada di tabel receipts
*/
func ParseReceipt(text string) (
	storeName string,
	total int64,
	date *time.Time,
	taxID string,
	isQualified bool,
	subtotal int64,
	tax int64,
	items []string,
) {
	lines := strings.Split(text, "\n")

	// 1️⃣ store_name: ambil baris yang mengandung huruf Jepang / "店" / "ショップ"
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Cari pola nama toko dengan karakter Jepang atau kata kunci
		if strings.Contains(line, "店") ||
			strings.Contains(line, "ショップ") ||
			regexp.MustCompile(`[\p{Hiragana}\p{Katakana}\p{Han}].+店`).MatchString(line) {
			storeName = line
			break
		}
	}
	fmt.Println("[DEBUG] Parsed store_name:", storeName)

	// 2️⃣ transaction_date: format YYYY年MM月DD
	reDate := regexp.MustCompile(`(\d{4})年\s*(\d{1,2})月\s*(\d{1,2})日?`)
	if m := reDate.FindStringSubmatch(text); len(m) > 3 {
		year, _ := strconv.Atoi(m[1])
		month, _ := strconv.Atoi(m[2])
		day, _ := strconv.Atoi(m[3])
		d := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
		date = &d
	}
	fmt.Println("[DEBUG] Parsed transaction_date:", date)

	// 3️⃣ total_amount: cari "合計" atau "総計" atau "合計金額"
	reTotal := regexp.MustCompile(`(合計|総計|合計金額)\s*[:：]?\s*¥?\s*([\d,]+)`)
	if m := reTotal.FindStringSubmatch(text); len(m) > 2 {
		totalStr := strings.ReplaceAll(m[2], ",", "")
		total, _ = strconv.ParseInt(totalStr, 10, 64)
	} else {
		// Fallback: cari angka terbesar yang mungkin total
		reAnyAmount := regexp.MustCompile(`¥?\s*([\d,]{3,})`)
		matches := reAnyAmount.FindAllStringSubmatch(text, -1)
		if len(matches) > 0 {
			// Ambil angka terbesar sebagai total
			var maxAmount int64
			for _, match := range matches {
				amountStr := strings.ReplaceAll(match[1], ",", "")
				amount, _ := strconv.ParseInt(amountStr, 10, 64)
				if amount > maxAmount {
					maxAmount = amount
				}
			}
			total = maxAmount
		}
	}
	fmt.Println("[DEBUG] Parsed total_amount:", total)

	// 4️⃣ subtotal (小計) & tax (税/ARE)
	reSubtotal := regexp.MustCompile(`(小計|税抜|税抜き)\s*[:：]?\s*¥?\s*([\d,]+)`)
	if m := reSubtotal.FindStringSubmatch(text); len(m) > 2 {
		subtotalStr := strings.ReplaceAll(m[2], ",", "")
		subtotal, _ = strconv.ParseInt(subtotalStr, 10, 64)
		fmt.Println("[DEBUG] Parsed subtotal:", subtotal)
	}

	reTax := regexp.MustCompile(`(消費税|税|TAX|ARE)\s*[:：]?\s*¥?\s*([\d,]+)`)
	if m := reTax.FindStringSubmatch(text); len(m) > 2 {
		taxStr := strings.ReplaceAll(m[2], ",", "")
		tax, _ = strconv.ParseInt(taxStr, 10, 64)
		fmt.Println("[DEBUG] Parsed tax:", tax)
	}

	// 5️⃣ tax_registration_id: cari "税番号", "登録番号", "T"
	reTaxID := regexp.MustCompile(`(税番号|登録番号|T\.)\s*[:：]?\s*([\w\d\-]+)`)
	if m := reTaxID.FindStringSubmatch(text); len(m) > 2 {
		taxID = strings.TrimSpace(m[2])
	}
	fmt.Println("[DEBUG] Parsed taxID:", taxID)

	// 6️⃣ is_qualified: cari "対象" untuk struk yang memenuhi syarat
	isQualified = strings.Contains(text, "対象") || strings.Contains(text, "対象外") == false
	fmt.Println("[DEBUG] Parsed isQualified:", isQualified)

	// 7️⃣ items: ambil semua baris yang mengandung harga (¥) dan nama item
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "¥") && !strings.Contains(line, "合計") {
			// Filter out total lines
			reItem := regexp.MustCompile(`(.+?)\s+¥?\s*([\d,]+)`)
			if match := reItem.FindStringSubmatch(line); len(match) > 2 {
				items = append(items, fmt.Sprintf("%s: ¥%s", match[1], match[2]))
			} else {
				items = append(items, line)
			}
		}
	}
	fmt.Println("[DEBUG] Parsed items count:", len(items))

	return
}

/*
ExtractText
- OCR via docker tesseract
*/
/*
ExtractText
- OCR via docker tesseract
*/
func ExtractText(imagePath string) (string, error) {
	fmt.Println("[DEBUG] === STARTING OCR EXTRACTION ===")
	fmt.Println("[DEBUG] Input image path:", imagePath)

	// Validasi file exists
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		fmt.Println("[ERROR] File does not exist:", imagePath)
		return "", fmt.Errorf("file not found: %s", imagePath)
	}

	filename := filepath.Base(imagePath)
	fmt.Println("[DEBUG] Filename:", filename)

	// path di dalam container
	containerPath := "/ocr/tmp/" + filename
	fmt.Println("[DEBUG] Container path:", containerPath)

	// Check if docker container is running
	fmt.Println("[DEBUG] Checking docker container...")
	checkCmd := exec.Command("docker", "ps", "--filter", "name=tesseract-ocr", "--format", "{{.Names}}")
	checkOut, checkErr := checkCmd.CombinedOutput()
	if checkErr != nil {
		fmt.Printf("[WARNING] Docker check failed: %v\n", checkErr)
	} else {
		containerStatus := strings.TrimSpace(string(checkOut))
		if containerStatus == "" {
			fmt.Println("[WARNING] Container 'tesseract-ocr' not running or not found")
		} else {
			fmt.Printf("[DEBUG] Container status: %s\n", containerStatus)
		}
	}

	// Build OCR command
	cmd := exec.Command(
		"docker", "exec", "tesseract-ocr",
		"tesseract",
		containerPath,
		"stdout",
		"-l", "jpn+eng",
		"--oem", "1", // gunakan LSTM engine
		"--psm", "6", // single uniform block of text
	)

	fmt.Println("[DEBUG] Executing command:", strings.Join(cmd.Args, " "))

	// Run OCR
	startTime := time.Now()
	out, err := cmd.CombinedOutput()
	elapsed := time.Since(startTime)

	fmt.Printf("[DEBUG] OCR completed in %v\n", elapsed)

	// Print raw output for debugging
	fmt.Println("[DEBUG] === RAW OCR OUTPUT ===")
	fmt.Println(string(out))
	fmt.Println("[DEBUG] === END RAW OUTPUT ===")

	if err != nil {
		errorMsg := fmt.Sprintf("docker OCR error: %v, output: %s", err, string(out))
		fmt.Println("[ERROR]", errorMsg)
		return "", fmt.Errorf("%s", errorMsg)
	}

	// Clean and validate output
	text := string(out)
	text = strings.TrimSpace(text)

	if text == "" {
		fmt.Println("[WARNING] OCR returned empty text")
	} else {
		fmt.Printf("[DEBUG] Extracted text length: %d characters\n", len(text))
		fmt.Printf("[DEBUG] First 200 chars: %.200s\n", text)
		if len(text) > 200 {
			fmt.Printf("[DEBUG] Last 200 chars: %.200s\n", text[len(text)-200:])
		}
	}

	return text, nil
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

/*
parseItemLine - parse string item menjadi description dan amount
Contoh: "コーヒー: ¥300" -> description="コーヒー", amount=300
*/
func parseItemLine(itemStr string) (description string, amount int64) {
	// Cari pola: description + ¥ + amount
	reItemDetail := regexp.MustCompile(`(.+?)\s*[:：]?\s*¥?\s*([\d,]+)`)
	if match := reItemDetail.FindStringSubmatch(itemStr); len(match) > 2 {
		description = strings.TrimSpace(match[1])
		amountStr := strings.ReplaceAll(match[2], ",", "")
		amount, _ = strconv.ParseInt(amountStr, 10, 64)
		return
	}

	// Fallback: ambil semua sebelum ¥ sebagai description
	if idx := strings.Index(itemStr, "¥"); idx > 0 {
		description = strings.TrimSpace(itemStr[:idx])
		amountPart := itemStr[idx:]
		reAmount := regexp.MustCompile(`[\d,]+`)
		if amountMatch := reAmount.FindString(amountPart); amountMatch != "" {
			amountStr := strings.ReplaceAll(amountMatch, ",", "")
			amount, _ = strconv.ParseInt(amountStr, 10, 64)
		}
		return
	}

	// Jika tidak ada ¥, gunakan seluruh string sebagai description
	description = itemStr
	return
}
