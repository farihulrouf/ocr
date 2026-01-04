package service

import (
	"ocr-saas-backend/internal/dto"
	"ocr-saas-backend/internal/repository"

	"github.com/google/uuid"
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
