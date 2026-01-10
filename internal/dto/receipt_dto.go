package dto

import (
	"github.com/google/uuid"
)

type MyReceiptRow struct {
	ID        uuid.UUID `json:"id"`
	RecordNo  string    `json:"record_no"`
	Date      string    `json:"date"`
	StoreName string    `json:"store_name"`
	Category  string    `json:"category"`
	Taxation  string    `json:"taxation"`
	Amount    int64     `json:"amount"`
	Status    string    `json:"status"`
}

type ReceiptUserInfo struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Name  string    `json:"name"`
}

type AdminReceiptRow struct {
	ID        uuid.UUID       `json:"id"`
	RecordNo  string          `json:"record_no"`
	Date      string          `json:"date"`
	StoreName string          `json:"store_name"`
	Category  string          `json:"category"`
	Taxation  string          `json:"taxation"`
	Amount    int64           `json:"amount"`
	Status    string          `json:"status"`
	User      ReceiptUserInfo `json:"user"`
}

/* =======================
   TAMBAHAN UNTUK DETAIL
   ======================= */

type ReceiptDetailResponse struct {
	ID        uuid.UUID `json:"id"`
	RecordNo  string    `json:"record_no"`
	Date      string    `json:"date"`
	StoreName string    `json:"store_name"`

	Category *ReceiptDetailCategory `json:"category"`
	Taxation string                 `json:"taxation"`
	Amount   int64                  `json:"amount"`
	Status   string                 `json:"status"`
	ImageURL string                 `json:"image_url"` // ✅ TAMBAH DI SINI
	Items    []ReceiptDetailItem    `json:"items"`

	User ReceiptUserInfo `json:"user"`

	OCRRaw any `json:"ocr_raw,omitempty"`
}

type ReceiptDetailCategory struct {
	ID   uuid.UUID `json:"id"`
	Code string    `json:"code"`
	Name string    `json:"name"`
}

type ReceiptDetailItem struct {
	ID          uint   `json:"id"` // ✅ WAJIB
	Description string `json:"description"`
	Amount      int64  `json:"amount"`
	TaxAmount   int64  `json:"tax_amount"`
	TaxRate     int    `json:"tax_rate"`
}

type BulkUpdateCategoryRequest struct {
	IDs   []uuid.UUID `json:"ids"`
	CatID uuid.UUID   `json:"cat_id"`
}

// dto/update_receipt_item.go
type UpdateReceiptItemRequest struct {
	Description *string `json:"description"`
	Amount      *int64  `json:"amount"`
	TaxAmount   *int64  `json:"tax_amount"`
	TaxRate     *int    `json:"tax_rate"`
}
