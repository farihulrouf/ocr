package ocr

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const OCRSpaceAPI = "https://api.ocr.space/parse/image"
const OCRSpaceAPIKey = "K82330939188957" // ganti dengan API key kamu

type OCRSpaceResponse struct {
	ParsedResults []struct {
		ParsedText string `json:"ParsedText"`
	} `json:"ParsedResults"`
	OCRExitCode  int         `json:"OCRExitCode"`
	IsErrored    bool        `json:"IsErroredOnProcessing"`
	ErrorMessage interface{} `json:"ErrorMessage"`
}

/*
// ExtractText melakukan OCR via OCR.Space API
func ExtractText(imagePath string) (string, error) {
	fmt.Println("[DEBUG] === STARTING OCR VIA OCR.SPACE API ===")
	fmt.Println("[DEBUG] Input image path:", imagePath)

	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found: %s", imagePath)
	}

	file, err := os.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to open image: %v", err)
	}
	defer file.Close()

	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	part, err := writer.CreateFormFile("file", imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %v", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return "", fmt.Errorf("failed to copy file: %v", err)
	}

	_ = writer.WriteField("language", "jpn")
	_ = writer.WriteField("OCREngine", "2")
	_ = writer.WriteField("isTable", "true")
	writer.Close()

	req, err := http.NewRequest("POST", OCRSpaceAPI, &b)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("apikey", OCRSpaceAPIKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 60 * time.Second}
	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call OCR API: %v", err)
	}
	defer resp.Body.Close()
	fmt.Printf("[DEBUG] OCR API request completed in %v\n", time.Since(start))

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read OCR API response: %v", err)
	}

	var ocrResp OCRSpaceResponse
	if err := json.Unmarshal(body, &ocrResp); err != nil {
		return "", fmt.Errorf("failed to parse OCR API response: %v", err)
	}

	if ocrResp.IsErrored || len(ocrResp.ParsedResults) == 0 {
		var errMsg string
		switch v := ocrResp.ErrorMessage.(type) {
		case string:
			errMsg = v
		case []interface{}:
			if len(v) > 0 {
				errMsg = fmt.Sprintf("%v", v[0])
			} else {
				errMsg = "unknown OCR error"
			}
		default:
			errMsg = "unknown OCR error"
		}
		return "", fmt.Errorf("OCR failed: %s", errMsg)
	}

	text := ocrResp.ParsedResults[0].ParsedText
	fmt.Println("[DEBUG] Extracted text length:", len(text))
	fmt.Println("[OCR RESULT START]\n", text, "\n[OCR RESULT END]")
	return text, nil
}
*/
// Struktur minimal untuk menangkap response teks dari Mistral
type MistralResponse struct {
	Pages []struct {
		Markdown string `json:"markdown"` // ここをcontentからmarkdownに修正
	} `json:"pages"`
}

func ExtractText(imagePath string) (string, error) {
	fmt.Println("\n[DEBUG] === STARTING OCR VIA MISTRAL API ===")

	// 1. Baca file
	fileData, err := os.ReadFile(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to read image file: %v", err)
	}

	// 2. Encode Base64
	base64Image := base64.StdEncoding.EncodeToString(fileData)
	dataURL := fmt.Sprintf("data:image/png;base64,%s", base64Image)

	// 3. Payload
	payload := map[string]interface{}{
		"model": "mistral-ocr-latest",
		"document": map[string]string{
			"type":      "image_url",
			"image_url": dataURL,
		},
	}
	jsonData, _ := json.Marshal(payload)

	// 4. Request
	req, err := http.NewRequest("POST", "https://api.mistral.ai/v1/ocr", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer KXoKnv3W0nq1kb2rVyrCvntVOvKpOZac")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error %d: %s", resp.StatusCode, string(body))
	}

	// 5. Parsing JSON Internal Mistral
	var result MistralResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse Mistral response: %v", err)
	}

	// 6. Ambil teks Markdown
	var finalFullText string
	for _, page := range result.Pages {
		// Ambil field Markdown, bukan Content
		finalFullText += page.Markdown + "\n"
	}

	if finalFullText == "" || finalFullText == "\n" {
		fmt.Println("[WARNING] No text extracted from markdown field")
	}

	fmt.Println("[DEBUG] OCR Text Extracted Successfully")
	return finalFullText, nil
}
