package service

import (
	"errors"
	"ocr-saas-backend/internal/dto"
	"ocr-saas-backend/internal/models"
	"ocr-saas-backend/internal/repository"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetMyReceipts(
	tenantID, userID uuid.UUID,
	page, pageSize int,
	q, status, sort string,
) (interface{}, error) {

	receipts, total, err := repository.GetMyReceipts(
		tenantID, userID, page, pageSize, q, status, sort,
	)
	if err != nil {
		return nil, err
	}

	rows := make([]dto.MyReceiptRow, 0, len(receipts))

	for _, r := range receipts {
		row := dto.MyReceiptRow{
			ID:        r.ID,
			RecordNo:  "R-" + r.ID.String()[:4],
			StoreName: r.StoreName,
			Amount:    r.TotalAmount,
			Status:    r.Status,
		}

		// date
		if r.TransactionDate != nil {
			row.Date = r.TransactionDate.Format("2006-01-02")
		}

		// category
		if r.AccountCategory != nil {
			row.Category = r.AccountCategory.Name
		}

		// taxation
		if r.IsQualified {
			row.Taxation = "eligible"
		} else {
			row.Taxation = "non-eligible"
		}

		rows = append(rows, row)
	}

	return map[string]interface{}{
		"data": rows,
		"meta": map[string]interface{}{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
		"status": "success",
	}, nil
}

func GetAllReceipts(
	tenantID uuid.UUID,
	page, pageSize int,
	q, status, sort string,
) (interface{}, error) {

	receipts, total, err := repository.GetAllReceipts(
		tenantID, page, pageSize, q, status, sort,
	)
	if err != nil {
		return nil, err
	}

	rows := make([]dto.AdminReceiptRow, 0, len(receipts))

	for _, r := range receipts {
		row := dto.AdminReceiptRow{
			ID:        r.ID,
			RecordNo:  "R-" + r.ID.String()[:4],
			StoreName: r.StoreName,
			Amount:    r.TotalAmount,
			Status:    r.Status,
		}

		if r.TransactionDate != nil {
			row.Date = r.TransactionDate.Format("2006-01-02")
		}

		if r.AccountCategory != nil {
			row.Category = r.AccountCategory.Name
		}

		if r.IsQualified {
			row.Taxation = "eligible"
		} else {
			row.Taxation = "non-eligible"
		}

		// 🔹 tambahan khusus admin
		if r.User.ID != uuid.Nil {
			row.User = dto.ReceiptUserInfo{
				ID:    r.User.ID,
				Email: r.User.Email,
				Name:  r.User.Name,
			}
		}

		rows = append(rows, row)
	}

	return map[string]interface{}{
		"data": rows,
		"meta": map[string]interface{}{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
		"status": "success",
	}, nil
}

func GetReceiptDetail(
	tenantID uuid.UUID,
	receiptID uuid.UUID,
) (*dto.ReceiptDetailResponse, error) {

	if tenantID == uuid.Nil || receiptID == uuid.Nil {
		return nil, errors.New("invalid id")
	}

	receipt, err := repository.GetReceiptDetailByID(tenantID, receiptID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("receipt not found")
		}
		return nil, err
	}

	response := mapReceiptToDetailDTO(receipt)
	return &response, nil
}

/* =====================
   PRIVATE MAPPER
   ===================== */

func mapReceiptToDetailDTO(r *models.Receipt) dto.ReceiptDetailResponse {
	// items
	items := make([]dto.ReceiptDetailItem, 0, len(r.LineItems))
	for _, it := range r.LineItems {
		items = append(items, dto.ReceiptDetailItem{
			Description: it.Description,
			Amount:      it.Amount,
			TaxAmount:   it.TaxAmount,
			TaxRate:     it.TaxRate,
		})
	}

	// category
	var category *dto.ReceiptDetailCategory
	if r.AccountCategory != nil {
		category = &dto.ReceiptDetailCategory{
			ID:   r.AccountCategory.ID,
			Code: r.AccountCategory.Code,
			Name: r.AccountCategory.Name,
		}
	}

	// date
	date := ""
	if r.TransactionDate != nil {
		date = r.TransactionDate.Format("2006-01-02")
	}

	return dto.ReceiptDetailResponse{
		ID:        r.ID,
		RecordNo:  "", // belum ada di model
		Date:      date,
		StoreName: r.StoreName,
		ImageURL:  r.ImageURL, // ✅ INI YANG KEMARIN HILAN
		Category:  category,
		Taxation:  mapTaxation(r.IsQualified),
		Amount:    r.TotalAmount,
		Status:    r.Status,
		Items:     items,
		User: dto.ReceiptUserInfo{
			ID:    r.User.ID,
			Email: r.User.Email,
			Name:  r.User.Name,
		},
	}
}

func mapTaxation(isQualified bool) string {
	if isQualified {
		return "eligible"
	}
	return "non-eligible"
}

var (
	ErrReceiptNotFound     = errors.New("receipt not found")
	ErrReceiptAlreadyFinal = errors.New("receipt already confirmed or rejected")
	ErrInvalidTotalAmount  = errors.New("total amount does not match receipt items")
)

func ConfirmReceipt(
	tenantID uuid.UUID,
	receiptID uuid.UUID,
	total int64,
	date time.Time,
) error {

	// 1️⃣ ambil detail receipt
	receipt, err := repository.GetReceiptDetailByID(tenantID, receiptID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrReceiptNotFound
		}
		return err
	}

	// 2️⃣ validasi status
	if receipt.Status != "PENDING" {
		return ErrReceiptAlreadyFinal
	}

	// 3️⃣ optional: validasi total dari items
	var sum int64
	for _, item := range receipt.LineItems {
		sum += item.Amount
	}

	if sum > 0 && sum != total {
		return ErrInvalidTotalAmount
	}
	// 4️⃣ update receipt
	return repository.ConfirmReceiptByID(
		tenantID,
		receiptID,
		total,
		date,
	)
}

func DeleteReceiptManager(
	tenantID uuid.UUID,
	receiptID uuid.UUID,
) error {

	err := repository.DeleteReceiptByIDManager(
		tenantID,
		receiptID,
	)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrReceiptNotFound
		}
		return err
	}

	return nil
}
