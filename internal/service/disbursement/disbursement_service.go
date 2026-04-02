package disbursement

import (
	"errors"
	"ocr-saas-backend/configs"
	"ocr-saas-backend/internal/models"
	"ocr-saas-backend/internal/repository/disbursement"
	"time"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

func ProcessPayment(reportID, financeID, tenantID uuid.UUID, amount int64, refNum, proofURL string) error {
	return configs.DB.Transaction(func(tx *gorm.DB) error {
		var report models.ExpenseReport

		// 1. Cek & Lock Data
		if err := tx.Set("gorm:query_option", "FOR UPDATE").
			Where("id = ? AND tenant_id = ?", reportID, tenantID).
			First(&report).Error; err != nil {
			return errors.New("laporan tidak ditemukan")
		}

		if report.Status != "APPROVED" {
			return errors.New("hanya laporan APPROVED yang bisa dibayar")
		}

		if report.TotalAmount != amount {
			return errors.New("nominal tidak sesuai")
		}

		// 2. Simpan Disbursement
		payout := models.Disbursement{
			TenantID:        tenantID,
			ExpenseReportID: &reportID,
			PayerID:         financeID,
			Amount:          report.TotalAmount,
			ReferenceNumber: refNum,
			ProofImageURL:   proofURL,
			PaidAt:          time.Now(),
		}
		if err := disbursement.Create(tx, &payout); err != nil {
			return err
		}

		// 3. Update Status
		if err := disbursement.UpdateReportStatus(tx, reportID, "PAID"); err != nil {
			return err
		}
		return disbursement.UpdateReceiptsStatus(tx, reportID, "PAID")
	})
}
