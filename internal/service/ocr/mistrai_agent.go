package ocr

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

/*
=========================
Domain Models
=========================
*/

// hasil final untuk ProcessOCR → saveReceiptItems
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

/*
=========================
PUBLIC API (AI)
=========================
*/

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

/*
=========================
DOMAIN ADAPTER
=========================
*/

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

/*
=========================
INTERNAL HELPERS
=========================
*/

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
