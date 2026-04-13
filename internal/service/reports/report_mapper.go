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
	// 1. Inisialisasi Approver sebagai NIL (kosong) secara default
	var approverPtr *dto.UserResponse

	// 2. CEK: Jika sudah di-approve/reject dan datanya ada
	if (r.Status == "APPROVED" || r.Status == "REJECTED" || r.Status == "PAID") && r.Approver != nil {
		// Gunakan tanda & untuk mengambil ALAMAT memori (Pointer)
		approverPtr = &dto.UserResponse{
			ID:     r.Approver.ID.String(),
			Name:   r.Approver.Name,
			Email:  r.Approver.Email,
			Role:   r.Approver.Role, // Masukkan Role & Avatar sekalian Mas
			Avatar: r.Approver.Avatar,
		}
	}

	return dto.ExpenseReportResponse{
		ID:          r.ID.String(),
		Title:       r.Title,
		TotalAmount: r.TotalAmount,
		Status:      r.Status,
		CreatedAt:   r.CreatedAt,
		User:        user,
		Approver:    approverPtr, // <--- Sekarang sudah cocok (Pointer ketemu Pointer)
		Receipts:    receipts,
	}
}
