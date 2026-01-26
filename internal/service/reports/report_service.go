package reports

import (
	"errors"
	"ocr-saas-backend/configs"
	"ocr-saas-backend/internal/dto"
	"ocr-saas-backend/internal/models"
	repo "ocr-saas-backend/internal/repository/reports"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

	if len(report.Receipts) == 0 {
		return errors.New("cannot submit empty report")
	}

	var total int64
	for _, r := range report.Receipts {
		total += r.TotalAmount
	}
	report.TotalAmount = total

	report.Status = "SUBMITTED"
	return repo.UpdateStatus(report.ID, report.Status, report.TotalAmount)
}

func UpdateReport(
	tenantID, userID, reportID uuid.UUID,
	title string,
) error {

	report, err := repo.GetByID(tenantID, reportID)
	if err != nil {
		return err
	}

	if report.UserID != userID {
		return errors.New("not owner of report")
	}

	if report.Status != "DRAFT" {
		return errors.New("only draft can be updated")
	}

	report.Title = title
	return repo.Update(report)
}

func GetPendingReports(
	tenantID uuid.UUID,
	page, pageSize int,
) ([]dto.ExpenseReportResponse, int64, error) {

	rows, total, err := repo.ListSubmitted(
		tenantID, page, pageSize,
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

func ApproveReport(
	tenantID, reportID uuid.UUID,
) error {

	report, err := repo.GetByID(tenantID, reportID)
	if err != nil {
		return err
	}

	if report.Status != "SUBMITTED" {
		return errors.New("report is not submitted")
	}
	return repo.UpdateReportStatus(
		tenantID,
		reportID,
		"APPROVED",
	)
}

func RejectReport(
	tenantID, reportID uuid.UUID,
) error {

	report, err := repo.GetByID(tenantID, reportID)
	if err != nil {
		return err
	}

	if report.Status != "SUBMITTED" {
		return errors.New("report is not submitted")
	}
	return repo.UpdateReportStatus(
		tenantID,
		reportID,
		"REJECTED",
	)
}

func GetMyReportDetail(
	tenantID, userID, reportID uuid.UUID,
) (*dto.ExpenseReportResponse, error) {

	report, err := repo.GetByID(tenantID, reportID)
	if err != nil {
		return nil, err
	}

	// ownership check
	if report.UserID != userID {
		return nil, errors.New("not owner of report")
	}

	res := ToExpenseReportResponse(*report)
	return &res, nil
}

func AddReceiptsToReport(tenantID, reportID uuid.UUID, receiptIDs []uuid.UUID) error {
	if len(receiptIDs) == 0 {
		return errors.New("no receipt ids provided")
	}

	// 1️⃣ Pastikan report itu milik tenant
	var report models.ExpenseReport
	if err := configs.DB.
		Where("id = ? AND tenant_id = ?", reportID, tenantID).
		First(&report).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("report not found for this tenant")
		}
		return err
	}

	// 2️⃣ Update setiap receipt -> set report_id
	for _, rid := range receiptIDs {
		res := configs.DB.Model(&models.Receipt{}).
			Where("id = ? AND tenant_id = ?", rid, tenantID).
			Updates(map[string]interface{}{
				"report_id": reportID,
			})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errors.New("one or more receipts not found or do not belong to tenant")
		}
	}

	return nil
}
