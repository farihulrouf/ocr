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
) (*models.ExpenseReport, error) {

	// GENERATE ID BARU DI SINI
	newReportID := uuid.New()

	report := &models.ExpenseReport{
		ID:       newReportID, // Isi ID-nya secara manual
		TenantID: tenantID,
		UserID:   userID,
		Title:    title,
		Status:   "DRAFT",
	}

	// Simpan ke database
	if err := repo.Create(report); err != nil {
		return nil, err
	}

	// Sekarang report.ID sudah terisi dengan newReportID
	return report, nil
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

	return configs.DB.Transaction(func(tx *gorm.DB) error {

		// 1️⃣ Pastikan report ada, milik tenant, dan MASIH DRAFT
		var report models.ExpenseReport
		if err := tx.Where("id = ? AND tenant_id = ? AND status = ?", reportID, tenantID, "DRAFT").
			First(&report).Error; err != nil {
			return errors.New("report not found or already submitted")
		}

		// 2️⃣ Update setiap receipt -> set report_id DAN status jadi 'REPORTED'
		for _, rid := range receiptIDs {
			// Gunakan Updates(map) agar lebih efisien untuk banyak kolom
			res := tx.Model(&models.Receipt{}).
				Where("id = ? AND tenant_id = ? AND status = ?", rid, tenantID, "SUCCESS").
				Updates(map[string]interface{}{
					"report_id": reportID,
					"status":    "REPORTED", // Otomatis jadi REPORTED
				})

			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected == 0 {
				// Ini mencegah struk yang sudah dipakai orang lain atau statusnya bukan SUCCESS untuk masuk
				return errors.New("receipt not found or already used: " + rid.String())
			}
		}

		// 3️⃣ HITUNG TOTAL BARU
		var newTotal int64
		err := tx.Model(&models.Receipt{}).
			Where("report_id = ? AND tenant_id = ?", reportID, tenantID).
			Select("COALESCE(SUM(total_amount), 0)").
			Scan(&newTotal).Error

		if err != nil {
			return err
		}

		// 4️⃣ UPDATE TOTAL_AMOUNT DI EXPENSE_REPORT
		// Gunakan Select("total_amount") agar GORM hanya update kolom itu saja
		if err := tx.Model(&report).Select("total_amount").Updates(models.ExpenseReport{TotalAmount: newTotal}).Error; err != nil {
			return err
		}

		return nil
	})
}

func RemoveReceiptFromReport(tenantID, reportID, receiptID uuid.UUID) error {
	return configs.DB.Transaction(func(tx *gorm.DB) error {

		// 1. Pastikan Report-nya ada dan statusnya masih DRAFT
		var report models.ExpenseReport
		if err := tx.Where("id = ? AND tenant_id = ? AND status = ?", reportID, tenantID, "DRAFT").
			First(&report).Error; err != nil {
			return errors.New("laporan tidak ditemukan atau sudah di-submit")
		}

		// 2. INI KUNCINYA: Set report_id jadi NULL dan status balik ke SUCCESS
		// Kita pakai map supaya gorm mau update ke NULL
		err := tx.Model(&models.Receipt{}).
			Where("id = ? AND report_id = ? AND tenant_id = ?", receiptID, reportID, tenantID).
			Updates(map[string]interface{}{
				"report_id": nil,       // Menghapus kaitan ke laporan
				"status":    "SUCCESS", // Mengembalikan status agar bisa dipakai lagi
			}).Error

		if err != nil {
			return err
		}

		// 3. HITUNG ULANG TOTAL_AMOUNT Laporan
		var newTotal int64
		tx.Model(&models.Receipt{}).
			Where("report_id = ? AND tenant_id = ?", reportID, tenantID).
			Select("COALESCE(SUM(total_amount), 0)").
			Scan(&newTotal)

		// 4. Update total_amount di tabel expense_reports
		return tx.Model(&report).Update("total_amount", newTotal).Error
	})
}
