package reports

import (
	"ocr-saas-backend/internal/dto"
	"ocr-saas-backend/internal/models"

	"github.com/google/uuid"
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

	// 🔐 SAFE user mapping (anti panic)
	user := dto.UserResponse{}
	if r.User.ID != uuid.Nil {
		user = dto.UserResponse{
			ID:    r.User.ID.String(),
			Name:  r.User.Name,
			Email: r.User.Email,
		}
	}
	approver := dto.UserResponse{}

	// 2. CEK: Jika status bukan DRAFT/SUBMITTED dan relasi Approver ada
	if (r.Status == "APPROVED" || r.Status == "REJECTED") && r.Approver != nil {
		approver = dto.UserResponse{
			ID:    r.Approver.ID.String(),
			Name:  r.Approver.Name,
			Email: r.Approver.Email,
		}
	}

	return dto.ExpenseReportResponse{
		ID:          r.ID.String(),
		Title:       r.Title,
		TotalAmount: r.TotalAmount,
		Status:      r.Status,
		CreatedAt:   r.CreatedAt,
		User:        user,     // ⬅️ TAMBAHAN
		Approver:    approver, // <--- TAMBAHKAN INI
		Receipts:    receipts,
	}
}
