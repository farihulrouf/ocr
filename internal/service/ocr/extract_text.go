package ocr

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Struktur response dari Mistral
type MistralResponse struct {
	Pages []struct {
		Markdown string `json:"markdown"`
	} `json:"pages"`
}

func ExtractText(imagePath string) (string, error) {

	fmt.Println("\n[DEBUG] === STARTING OCR ===")

	apiKey := os.Getenv("MISTRAL_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("MISTRAL_API_KEY not found in environment")
	}

	fileData, err := os.ReadFile(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to read image file: %v", err)
	}

	base64Image := base64.StdEncoding.EncodeToString(fileData)
	dataURL := fmt.Sprintf("data:image/png;base64,%s", base64Image)

	payload := map[string]interface{}{
		"model": "mistral-ocr-latest",
		"document": map[string]string{
			"type":      "image_url",
			"image_url": dataURL,
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	maxRetry := 3

	for attempt := 1; attempt <= maxRetry; attempt++ {

		fmt.Printf("[DEBUG] OCR Attempt %d\n", attempt)

		start := time.Now()

		req, err := http.NewRequest(
			"POST",
			"https://api.mistral.ai/v1/ocr",
			bytes.NewBuffer(jsonData),
		)

		if err != nil {
			return "", err
		}

		req.Header.Set("Authorization", "Bearer "+apiKey)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)

		duration := time.Since(start)
		fmt.Println("[DEBUG] OCR request duration:", duration)

		if err != nil {

			fmt.Println("[ERROR] Request failed:", err)

			if attempt == maxRetry {
				return "", err
			}

			time.Sleep(2 * time.Second)
			continue
		}

		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		if resp.StatusCode != http.StatusOK {

			fmt.Printf("[ERROR] Status %d: %s\n", resp.StatusCode, string(body))

			if attempt == maxRetry {
				return "", fmt.Errorf("error %d: %s", resp.StatusCode, string(body))
			}

			time.Sleep(2 * time.Second)
			continue
		}

		var result MistralResponse

		if err := json.Unmarshal(body, &result); err != nil {
			return "", fmt.Errorf("failed to parse Mistral response: %v", err)
		}

		jsonBytes, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println("[DEBUG] JSON MistralResponse:\n", string(jsonBytes))

		var finalFullText string

		for _, page := range result.Pages {
			finalFullText += page.Markdown + "\n"
		}

		if finalFullText == "" {
			fmt.Println("[WARNING] No text extracted from markdown field")
		}

		fmt.Println("[DEBUG] OCR Text Extracted Successfully")

		return finalFullText, nil
	}

	return "", fmt.Errorf("OCR failed after retries")
}
