package reports

import (
	"ocr-saas-backend/configs"
	"ocr-saas-backend/internal/models"

	"github.com/google/uuid"
)

// CreateExportLog - Simpan riwayat request export ke tabel export_logs
func CreateExportLog(tenantID, userID uuid.UUID, format string) (*models.ExportLog, error) {
	log := models.ExportLog{
		Base: models.Base{
			ID: uuid.New(),
		},
		TenantID: tenantID,
		UserID:   userID,
		Format:   format,
		FileURL:  "", // Masih kosong, nanti diupdate Lambda
	}

	err := configs.DB.Create(&log).Error
	if err != nil {
		return nil, err
	}

	return &log, nil
}

func UpdateExportLogURL(id uuid.UUID, fileURL string) error {
	return configs.DB.Model(&models.ExportLog{}).
		Where("id = ?", id).
		Update("file_url", fileURL).
		Error
}
