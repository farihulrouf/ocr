package reports

import (
	"ocr-saas-backend/internal/dto"
	"ocr-saas-backend/internal/models"
)

func ToExpenseReportResponse(r models.ExpenseReport) dto.ExpenseReportResponse {
	receipts := make([]dto.ReceiptResponse, 0)

	for _, rc := range r.Receipts {
		receipts = append(receipts, dto.ReceiptResponse{
			ID:          rc.ID.String(),
			StoreName:   rc.StoreName,
			TotalAmount: rc.TotalAmount,
			Status:      rc.Status,
		})
	}

	return dto.ExpenseReportResponse{
		ID:          r.ID.String(),
		Title:       r.Title,
		TotalAmount: r.TotalAmount,
		Status:      r.Status,
		CreatedAt:   r.CreatedAt,
		Receipts:    receipts,
	}
}
