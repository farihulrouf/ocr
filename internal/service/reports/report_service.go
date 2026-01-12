package reports

import (
	"errors"
	"ocr-saas-backend/internal/dto"
	"ocr-saas-backend/internal/models"
	repo "ocr-saas-backend/internal/repository/reports"

	"github.com/google/uuid"
)

func GetMyReports(
	tenantID, userID uuid.UUID,
	page, pageSize int,
) ([]dto.ExpenseReportResponse, int64, error) {

	rows, total, err := repo.ListMyReports(
		tenantID, userID, page, pageSize,
	)
	if err != nil {
		return nil, 0, err
	}

	result := make([]dto.ExpenseReportResponse, 0)
	for _, r := range rows {
		result = append(result, ToExpenseReportResponse(r))
	}

	return result, total, nil
}

func CreateReport(
	tenantID, userID uuid.UUID,
	title string,
) error {

	report := &models.ExpenseReport{
		TenantID: tenantID,
		UserID:   userID,
		Title:    title,
		Status:   "DRAFT",
	}

	return repo.Create(report)
}

func SubmitReport(
	tenantID, userID, reportID uuid.UUID,
) error {

	report, err := repo.GetByID(tenantID, reportID)
	if err != nil {
		return err
	}

	if report.UserID != userID {
		return errors.New("not owner of report")
	}

	if report.Status != "DRAFT" {
		return errors.New("report already submitted")
	}

	report.Status = "SUBMITTED"
	return repo.Update(report)
}
