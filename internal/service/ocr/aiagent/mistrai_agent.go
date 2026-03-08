package aiagent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

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

	fmt.Println("\n[DEBUG] >>> STEP 2: AI STRUCTURING <<<")
	fmt.Printf("[DEBUG] Raw Text Length: %d characters\n", len(ocrText))

	apiKey := os.Getenv("MISTRAL_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("MISTRAL_API_KEY not set")
	}

	prompt := fmt.Sprintf(`Task: Parse this raw Japanese receipt text into a precise JSON format.

STRICT RULES:
1. EXTRACT: Store Name, Store ID (Tax ID starting with 'T'), Date (YYYY-MM-DD), Items, Total, and Tax.
2. ITEM FILTER: Only include actual products in the "items" array.
3. CLEANING: Fix OCR typos.
4. MATH: Prices must be numeric integers.

JSON STRUCTURE:
{
 "store_info": {"store_name": "...", "store_id": "..."},
 "transaction_info": {"date": "YYYY-MM-DD"},
 "items": [{"name": "...", "price": 0}],
 "payment_summary": {"total_amount": 0, "tax_details": {"tax_value": 0}}
}

RAW TEXT:
%s`, ocrText)

	payload := map[string]interface{}{
		"model": "mistral-small-latest",
		"messages": []map[string]string{
			{"role": "system", "content": "You are a professional financial data auditor. Return ONLY valid JSON."},
			{"role": "user", "content": prompt},
		},
		"response_format": map[string]string{"type": "json_object"},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	var lastErr error

	// 🔁 Retry maksimal 3x
	for attempt := 1; attempt <= 3; attempt++ {

		fmt.Printf("[DEBUG] AI Attempt %d\n", attempt)

		start := time.Now()

		req, err := http.NewRequest(
			"POST",
			"https://api.mistral.ai/v1/chat/completions",
			bytes.NewBuffer(jsonData),
		)
		if err != nil {
			return "", err
		}

		req.Header.Set("Authorization", "Bearer "+apiKey)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)

		if err != nil {
			fmt.Printf("[ERROR] AI Request Failed: %v\n", err)
			lastErr = err
			time.Sleep(2 * time.Second)
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		duration := time.Since(start)
		fmt.Println("[DEBUG] AI request duration:", duration)

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("[ERROR] AI API error %d: %s\n", resp.StatusCode, string(body))
			lastErr = fmt.Errorf("error %d", resp.StatusCode)
			time.Sleep(2 * time.Second)
			continue
		}

		var result struct {
			Choices []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			} `json:"choices"`
		}

		if err := json.Unmarshal(body, &result); err != nil {
			return "", fmt.Errorf("failed to parse AI response: %v", err)
		}

		if len(result.Choices) > 0 {
			fmt.Println("[DEBUG] AI JSON Response received successfully.")
			return result.Choices[0].Message.Content, nil
		}

		lastErr = fmt.Errorf("AI returned empty content")
		time.Sleep(2 * time.Second)
	}

	return "", fmt.Errorf("AI failed after retries: %v", lastErr)
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
