package reports

import (
	"context"
	"fmt"
	"log"

	"ocr-saas-backend/internal/infrastructure/aws"
	"ocr-saas-backend/internal/models"
	repo "ocr-saas-backend/internal/repository/reports"

	"github.com/google/uuid"
)

// TriggerExpenseExport - Memulai proses export data expense
func TriggerExpenseExport(ctx context.Context, tenantID, userID uuid.UUID, status string) error {
	log.Printf("[Service] Memulai proses export untuk Tenant: %s, User: %s, Status: %s", tenantID, userID, status)

	// 1. Buat log awal di database dengan format EXCEL
	logEntry, err := repo.CreateExportLog(tenantID, userID, "EXCEL")
	if err != nil {
		return fmt.Errorf("gagal membuat export log di database: %v", err)
	}

	// 2. Siapkan Payload untuk dikirim ke AWS Lambda
	payload := aws.ExportPayload{
		ExportLogID: logEntry.ID.String(),
		TenantID:    tenantID.String(),
		UserID:      userID.String(),
		Status:      status,
	}

	// 3. Panggil Lambda secara Async (Fire and Forget)
	err = aws.InvokeExportLambda(ctx, payload)
	if err != nil {
		log.Printf("[Service] ❌ Gagal memicu Lambda: %v", err)
		return fmt.Errorf("gagal memicu worker export: %v", err)
	}

	log.Printf("[Service] ✅ Sukses memicu Lambda. LogID: %s", logEntry.ID.String())
	return nil
}

// GetRecentExportLogs - Mengambil daftar riwayat export terbaru
func GetRecentExportLogs(ctx context.Context, tenantID uuid.UUID, limit int) ([]models.ExportLog, error) {
	log.Printf("[Service] Mengambil %d riwayat export terakhir untuk Tenant: %s", limit, tenantID)
	return repo.GetExportLogs(tenantID, limit)
}
