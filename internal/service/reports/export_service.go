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

// TriggerExpenseExport - Logika bisnis untuk memulai proses export
// Kita tambahkan parameter ctx agar bisa meneruskan context dari Handler Fiber
func TriggerExpenseExport(ctx context.Context, tenantID, userID uuid.UUID, status string) error {

	log.Printf("[Service] Memulai proses export untuk Tenant: %s, User: %s", tenantID, userID)

	// 1. Buat log di DB (Status awal biasanya kosong/PENDING di database)
	// Fungsi ini memastikan kita punya ID Record sebelum Lambda jalan
	logEntry, err := repo.CreateExportLog(tenantID, userID, "EXCEL")
	if err != nil {
		return fmt.Errorf("gagal membuat export log di database: %v", err)
	}

	// 2. Siapkan Payload menggunakan struct yang sudah kita buat di package aws
	// Ini lebih aman daripada pakai map[string]interface{} agar tidak typo key-nya
	payload := aws.ExportPayload{
		ExportLogID: logEntry.ID.String(),
		TenantID:    tenantID.String(),
		UserID:      userID.String(),
		Status:      status,
	}

	// 3. Panggil Lambda secara Async
	// Kita panggil fungsi invoke yang tadi kita buat di internal/infrastructure/aws
	err = aws.InvokeExportLambda(ctx, payload)
	if err != nil {
		// Jika gagal panggil Lambda, kita log error-nya
		// Tapi log di DB sudah terlanjur dibuat (bisa buat cleanup logic di sini jika perlu)
		return fmt.Errorf("gagal memicu worker export: %v", err)
	}

	log.Printf("[Service] Sukses memicu Lambda untuk LogID: %s", logEntry.ID.String())

	return nil
}

func GetRecentExportLogs(ctx context.Context, tenantID uuid.UUID, limit int) ([]models.ExportLog, error) {
	// Service tinggal memanggil repository
	return repo.GetExportLogs(tenantID, limit)
}
