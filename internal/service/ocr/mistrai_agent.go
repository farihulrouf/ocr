package ocr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

/*
=========================
Domain Models
=========================
*/

// hasil final untuk ProcessOCR → saveReceiptItems
/*
type ParsedItem struct {
	Description string
	Amount      int64
}

// raw hasil dari AI
type ParsedReceipt struct {
	StoreName         string  `json:"store_name"`
	TransactionDate   string  `json:"transaction_date"`
	TotalAmount       int     `json:"total_amount"`
	Subtotal          int     `json:"subtotal"`
	Tax               int     `json:"tax"`
	TaxRegistrationID *string `json:"tax_registration_id"`
	IsQualified       bool    `json:"is_qualified"`
	Items             []Item  `json:"items"`
}

type Item struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}
*/

/*
=========================
PUBLIC API (AI)
=========================
*/
/*
func ParseReceiptAi(ocrText string) (*ParsedReceipt, error) {
	fmt.Println("[DEBUG][AI] ParseReceiptAi called")

	apiKey := os.Getenv("MISTRAL_API_KEY")
	if apiKey == "" {
		return nil, errors.New("MISTRAL_API_KEY not set")
	}
	fmt.Println("[DEBUG][AI] API key loaded, len =", len(apiKey))

	payload, err := buildPayload(ocrText)
	if err != nil {
		return nil, err
	}
	fmt.Println("[DEBUG][AI] Payload size:", len(payload))

	req, err := http.NewRequest(
		"POST",
		"https://api.mistral.ai/v1/chat/completions",
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Println("[DEBUG][AI] HTTP status:", resp.Status)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("mistral api error %d: %s", resp.StatusCode, body)
	}

	return parseResponse(resp.Body)
}

func ParseReceipt(text string) (
	store string,
	total int64,
	date *time.Time,
	taxID string,
	isQualified bool,
	subtotal int64,
	tax int64,
	items []ParsedItem,
) {
	fmt.Println("[DEBUG][OCR] ParseReceipt called")
	fmt.Println("[DEBUG][OCR] OCR text length:", len(text))

	parsed, err := ParseReceiptAi(text)
	if err != nil {
		fmt.Println("[ERROR][OCR] ParseReceiptAi failed:", err)
		return
	}

	fmt.Printf("[DEBUG][OCR] ParsedReceipt: %+v\n", *parsed)

	// ===== receipt =====
	store = parsed.StoreName
	total = int64(parsed.TotalAmount)
	subtotal = int64(parsed.Subtotal)
	tax = int64(parsed.Tax)
	isQualified = parsed.IsQualified

	if parsed.TransactionDate != "" {
		if t, err := time.Parse("2006-01-02T15:04", parsed.TransactionDate); err == nil {
			date = &t
		} else {
			fmt.Println("[WARN][OCR] invalid date:", parsed.TransactionDate)
		}
	}

	if parsed.TaxRegistrationID != nil {
		taxID = *parsed.TaxRegistrationID
	}

	// ===== items (SOURCE OF TRUTH = AI) =====
	for i, it := range parsed.Items {
		fmt.Printf(
			"[DEBUG][ITEM] #%d name=%q price=%d\n",
			i+1, it.Name, it.Price,
		)

		if it.Name == "" || it.Price <= 0 {
			fmt.Println("[WARN][ITEM] skipped invalid item")
			continue
		}

		items = append(items, ParsedItem{
			Description: it.Name,
			Amount:      int64(it.Price),
		})
	}

	fmt.Println("[DEBUG][ITEM] total items parsed:", len(items))
	return
}

func buildPayload(text string) ([]byte, error) {
	prompt := fmt.Sprintf(`
You are a receipt parsing engine.

Rules:
- Output ONLY valid JSON
- No markdown
- No explanation
- Currency is JPY (integer)
- Date format: YYYY-MM-DDTHH:MM
- Missing values must be null

Fields:
store_name
transaction_date
total_amount
subtotal
tax
tax_registration_id
is_qualified
items [{name, price}]

OCR TEXT:
%s
`, text)

	payload := map[string]interface{}{
		"model":       "mistral-large-latest",
		"temperature": 0,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You extract structured data from Japanese receipts.",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	return json.Marshal(payload)
}

func parseResponse(body io.Reader) (*ParsedReceipt, error) {
	rawBytes, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	fmt.Println("[DEBUG][AI] Raw response:")
	fmt.Println(string(rawBytes))

	var raw struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(rawBytes, &raw); err != nil {
		return nil, err
	}

	if len(raw.Choices) == 0 {
		return nil, errors.New("empty AI response")
	}

	jsonText := strings.TrimSpace(raw.Choices[0].Message.Content)

	// 🔹 STRIP MARKDOWN / BACKTICKS
	jsonText = strings.TrimPrefix(jsonText, "```json")
	jsonText = strings.TrimPrefix(jsonText, "```")
	jsonText = strings.TrimSuffix(jsonText, "```")
	jsonText = strings.TrimSpace(jsonText)

	fmt.Println("[DEBUG][AI] JSON content after strip:")
	fmt.Println(jsonText)

	var parsed ParsedReceipt
	if err := json.Unmarshal([]byte(jsonText), &parsed); err != nil {
		return nil, fmt.Errorf("invalid AI JSON: %w\nRAW:\n%s", err, jsonText)
	}

	if parsed.StoreName == "" {
		return nil, errors.New("store_name missing")
	}
	if parsed.TotalAmount <= 0 {
		return nil, errors.New("total_amount invalid")
	}

	return &parsed, nil
}
*/
type MistralParsedJSON struct {
	StoreInfo struct {
		StoreName string `json:"store_name"`
		StoreID   string `json:"store_id"` // Untuk Tax ID (T+13 digit)
	} `json:"store_info"`
	TransactionInfo struct {
		Date string `json:"date"`
	} `json:"transaction_info"`
	Items []struct {
		Name  string `json:"name"`
		Price int64  `json:"price"`
	} `json:"items"`
	PaymentSummary struct {
		TotalAmount int64 `json:"total_amount"`
		TaxDetails  struct {
			TaxValue int64 `json:"tax_value"`
		} `json:"tax_details"`
	} `json:"payment_summary"`
}

type ParsedItem struct {
	Description string
	Amount      int64
}

func StructureTextWithAI(ocrText string) (string, error) {
	fmt.Println("\n[DEBUG] >>> STEP 2: AI STRUCTURING (INTELLIGENCE) <<<")
	fmt.Printf("[DEBUG] Raw Text Length: %d characters\n", len(ocrText))

	prompt := fmt.Sprintf(`Task: Parse this raw Japanese receipt text into a precise JSON format.
    
    STRICT RULES:
    1. EXTRACT: Store Name, Store ID (Tax ID starting with 'T'), Date (YYYY-MM-DD), Items, Total, and Tax.
    2. ITEM FILTER: Only include actual products in the "items" array. 
       - EXCLUDE: "Quantity Total", "Subtotal", "Tax Amount", or "Cash Received".
    3. CLEANING: Fix OCR typos (e.g., "ãƒšã‚·ãƒ‰" to "ãƒƒãƒˆãƒœãƒˆãƒ«").
    4. MATH: Prices must be numeric integers. (Sum of items) + Tax must equal Total Amount.

    JSON STRUCTURE:
    {
      "store_info": {"store_name": "...", "store_id": "..."},
      "transaction_info": {"date": "YYYY-MM-DD"},
      "items": [{"name": "...", "price": 0}],
      "payment_summary": {"total_amount": 0, "tax_details": {"tax_value": 0}}
    }

    RAW TEXT TO ANALYZE:
    %s`, ocrText)

	payload := map[string]interface{}{
		"model": "mistral-small-latest",
		"messages": []map[string]string{
			{"role": "system", "content": "You are a professional financial data auditor. Return ONLY valid JSON."},
			{"role": "user", "content": prompt},
		},
		"response_format": map[string]string{"type": "json_object"},
	}

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "https://api.mistral.ai/v1/chat/completions", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer KXoKnv3W0nq1kb2rVyrCvntVOvKpOZac")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 40 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[ERROR] AI Request Failed: %v\n", err)
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	json.Unmarshal(body, &result)

	if len(result.Choices) > 0 {
		fmt.Println("[DEBUG] AI JSON Response received successfully.")
		return result.Choices[0].Message.Content, nil
	}
	return "", fmt.Errorf("AI returned empty content")
}

func ParseReceipt(jsonText string) (
	store string, total int64, date time.Time, taxID string,
	isQualified bool, subtotal int64, tax int64, items []ParsedItem,
) {
	fmt.Println("\n[DEBUG] >>> STEP 3: FINAL MAPPING & DB VALIDATION <<<")

	// Membersihkan format JSON dari AI
	cleanJSON := strings.TrimSpace(jsonText)
	cleanJSON = strings.TrimPrefix(cleanJSON, "```json")
	cleanJSON = strings.TrimSuffix(cleanJSON, "```")

	var data MistralParsedJSON
	if err := json.Unmarshal([]byte(cleanJSON), &data); err != nil {
		fmt.Printf("[ERROR] Failed to unmarshal JSON: %v\n", err)
		return "Unknown", 0, time.Now(), "", false, 0, 0, nil
	}

	// 1. Basic Mapping
	store = data.StoreInfo.StoreName
	taxID = data.StoreInfo.StoreID
	total = data.PaymentSummary.TotalAmount
	tax = data.PaymentSummary.TaxDetails.TaxValue

	// Cek apakah ini Qualified Invoice (Jepang)
	isQualified = (taxID != "") && strings.HasPrefix(taxID, "T")

	// 2. Item Processing & Filtering
	var sumOfItems int64
	fmt.Println("[DEBUG] Filtering & Summing Items:")
	for _, it := range data.Items {
		// Filter kata kunci yang sering nyasar jadi item
		lowerName := strings.ToLower(it.Name)
		if strings.Contains(lowerName, "tax") || strings.Contains(lowerName, "total") ||
			strings.Contains(lowerName, "内税") || strings.Contains(lowerName, "合計") {
			fmt.Printf("   [SKIP] Non-product detected: %s\n", it.Name)
			continue
		}

		items = append(items, ParsedItem{
			Description: it.Name,
			Amount:      it.Price,
		})
		sumOfItems += it.Price
		fmt.Printf("   [ITEM] %-25s | Price: %d\n", it.Name, it.Price)
	}

	// 3. Mathematical Validation
	// Di struk Jepang, 'Subtotal' biasanya (Total - Tax) atau (Total) jika pajak sudah termasuk (内税)
	// Kita asumsikan subtotal murni adalah total dikurangi pajak
	expectedSubtotal := total - tax

	fmt.Printf("[DEBUG] Math Check -> Item Sum: %d | Expected Subtotal: %d\n", sumOfItems, expectedSubtotal)

	if sumOfItems != expectedSubtotal {
		diff := expectedSubtotal - sumOfItems
		fmt.Printf("[WARNING] Math Mismatch! Diff: %d Yen. Adding adjustment row.\n", diff)

		if diff != 0 {
			items = append(items, ParsedItem{
				Description: "Adjustment (Rounding/Others)",
				Amount:      diff,
			})
			fmt.Println("   [DEBUG] Adjustment row added to maintain balance.")
		}
	}

	// 4. Date Parsing
	parsedDate, err := time.Parse("2006-01-02", data.TransactionInfo.Date)
	if err != nil {
		fmt.Printf("[WARNING] Date format error '%s'. Using current time.\n", data.TransactionInfo.Date)
		date = time.Now()
	} else {
		date = parsedDate
	}

	subtotal = expectedSubtotal
	fmt.Printf("[DEBUG] Final Result: %s | Total: %d | Items: %d\n", store, total, len(items))
	return
}
