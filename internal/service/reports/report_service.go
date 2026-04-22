package reports

import (
	"errors"
	"fmt"
	"log"
	"ocr-saas-backend/configs"
	"ocr-saas-backend/internal/dto"
	"ocr-saas-backend/internal/models"
	repo "ocr-saas-backend/internal/repository/reports"
	"ocr-saas-backend/internal/service/budgets"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetMyReports(
	tenantID, userID uuid.UUID,
	page, pageSize int,
	status string,
	role string, // <--- Terima Role
) ([]dto.ExpenseReportResponse, int64, error) {

	// LOGIKA: Jika dia bukan EMPLOYEE, kirim uuid.Nil ke repo
	// supaya repo tidak memfilter berdasarkan user_id.
	filterUserID := userID
	if role != "EMPLOYEE" {
		filterUserID = uuid.Nil
	}

	// Panggil Repo (Cukup 5 Argumen, role tidak perlu dikirim ke DB)
	rows, total, err := repo.ListMyReports(
		tenantID, filterUserID, page, pageSize, status,
	)

	if err != nil {
		return nil, 0, err
	}

	// Mapping ke DTO
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

func ApproveReport(tenantID, reportID, managerID uuid.UUID) error {
	return configs.DB.Transaction(func(tx *gorm.DB) error {
		// 0. Ambil data laporan dulu untuk tahu nominal pengeluarannya
		var report models.ExpenseReport
		if err := tx.Select("total_amount").
			Where("id = ? AND tenant_id = ?", reportID, tenantID).
			First(&report).Error; err != nil {
			return err
		}

		// 1. POTONG BUDGET (Logic Baru)
		// Gunakan fungsi ConsumeBudgetLogic agar spent_amount di tabel budgets bertambah
		if err := budgets.ConsumeBudgetLogic(tx, tenantID, report.TotalAmount); err != nil {
			return err // Jika budget tidak cukup, transaksi gagal (Rollback)
		}

		err := repo.UpdateReportStatus(
			tenantID,
			reportID,
			"APPROVED",
			&managerID, // <-- Kirim pointer managerID ke repo
		)
		if err != nil {
			return err
		}

		// 3. Catat ke ApprovalLog (Sejarah)
		approvalLog := models.ApprovalLog{
			ExpenseReportID: &reportID,
			UserID:          managerID,
			Action:          "APPROVE",
			Comment:         "Laporan disetujui oleh Manajer",
		}
		return tx.Create(&approvalLog).Error
	})
}

func RejectReport(
	tenantID, reportID, managerID uuid.UUID, // <-- Tambah managerID di sini
) error {

	report, err := repo.GetByID(tenantID, reportID)
	if err != nil {
		return err
	}

	if report.Status != "SUBMITTED" {
		return errors.New("laporan belum di-submit, tidak bisa ditolak")
	}

	// Gunakan fungsi repo yang bisa mengupdate status sekaligus mencatat pelakunya
	return repo.UpdateReportStatus(
		tenantID,
		reportID,
		"REJECTED",
		&managerID, // <-- Kirim pointer managerID ke repo
	)
}

// Tambahkan parameter role string
func GetMyReportDetail(
	tenantID, userID, reportID uuid.UUID, role string,
) (*dto.ExpenseReportResponse, error) {

	report, err := repo.GetByID(tenantID, reportID)
	if err != nil {
		return nil, err
	}

	// --- LOGIKA BARU ---
	// Jika dia adalah EMPLOYEE, dia WAJIB menjadi pemilik report (UserID harus sama)
	// Jika dia MANAGER, FINANCE, atau ADMIN, dia BOLEH melihat asalkan TenantID sama (TenantID sudah dicek di repo.GetByID)
	if role == "EMPLOYEE" && report.UserID != userID {
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

		// 1️⃣ Pastikan report ada
		var report models.ExpenseReport
		if err := tx.Where("id = ? AND tenant_id = ?", reportID, tenantID).
			First(&report).Error; err != nil {
			return errors.New("report not found")
		}

		// 2️⃣ Update setiap receipt: set report_id DAN status jadi 'REPORTED'
		for _, rid := range receiptIDs {
			res := tx.Model(&models.Receipt{}).
				Where("id = ? AND tenant_id = ?", rid, tenantID).
				Updates(map[string]interface{}{
					"report_id": reportID,
					"status":    "REPORTED", // ✅ Status berubah jadi REPORTED
				})

			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected == 0 {
				return errors.New("receipt not found: " + rid.String())
			}
		}

		// 3️⃣ Hitung Total Baru dari semua receipt yang ada di report ini
		var newTotal int64
		err := tx.Model(&models.Receipt{}).
			Where("report_id = ? AND tenant_id = ?", reportID, tenantID).
			Select("COALESCE(SUM(total_amount), 0)").
			Scan(&newTotal).Error

		if err != nil {
			return err
		}

		// 4️⃣ Update Total Amount di Expense Report
		// Menggunakan Select("total_amount") agar GORM hanya update kolom itu saja
		if err := tx.Model(&report).Update("total_amount", newTotal).Error; err != nil {
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

// GetReadyToPayReports mengambil semua laporan yang sudah di-approve Manager dan siap dibayar
// Ubah: Tambahkan parameter 'status'
// internal/service/report_service.go (atau file tempat fungsi ini berada)

func GetReadyToPayReports(tenantID uuid.UUID, page, pageSize int, status string) ([]dto.ExpenseReportResponse, int64, error) {
	// 1. Ambil data mentah dari repo
	rows, total, err := repo.ListByStatus(tenantID, page, pageSize, status)
	if err != nil {
		return nil, 0, err
	}

	// 2. Transformasi ke DTO
	var reportsDTO []dto.ExpenseReportResponse
	for _, r := range rows {
		reportsDTO = append(reportsDTO, ToExpenseReportResponse(r))
	}

	return reportsDTO, total, nil
}

func BulkApproveReports(
	tenantID uuid.UUID,
	reportIDs []uuid.UUID,
	managerID uuid.UUID,
) error {

	if len(reportIDs) == 0 {
		return errors.New("no report ids provided")
	}

	// 🔥 VALIDASI UUID (anti uuid.Nil)
	for _, id := range reportIDs {
		if id == uuid.Nil {
			return errors.New("invalid report id detected (uuid.Nil)")
		}
	}

	return configs.DB.Transaction(func(tx *gorm.DB) error {

		// 🔍 DEBUG LOG
		log.Println("=== BULK APPROVE DEBUG ===")
		log.Println("TenantID:", tenantID)
		log.Println("ManagerID:", managerID)
		log.Println("ReportIDs:", reportIDs)

		// 1️⃣ Ambil semua report TANPA FILTER dulu (biar tahu mana yang ada)
		var reports []models.ExpenseReport
		if err := tx.
			Where("id IN ?", reportIDs).
			Find(&reports).Error; err != nil {
			return err
		}

		log.Println("Found Reports:", len(reports))

		// ❌ Kalau tidak ada sama sekali
		if len(reports) == 0 {
			return errors.New("no reports found for given IDs")
		}

		// 🔎 Validasi satu per satu (lebih jelas dari len mismatch)
		var validReports []models.ExpenseReport
		var invalidIDs []uuid.UUID

		for _, r := range reports {

			// cek tenant
			if r.TenantID != tenantID {
				log.Println("Tenant mismatch:", r.ID)
				invalidIDs = append(invalidIDs, r.ID)
				continue
			}

			// cek status
			if r.Status != "SUBMITTED" {
				log.Println("Invalid status:", r.ID, r.Status)
				invalidIDs = append(invalidIDs, r.ID)
				continue
			}

			validReports = append(validReports, r)
		}

		// ❌ kalau ada yang invalid → kasih error jelas
		if len(invalidIDs) > 0 {
			return fmt.Errorf("some reports invalid (tenant/status): %v", invalidIDs)
		}

		// ❌ kalau jumlah tidak match (ID tidak ditemukan)
		if len(validReports) != len(reportIDs) {
			return fmt.Errorf(
				"some report_ids not found (expected %d, got %d)",
				len(reportIDs),
				len(validReports),
			)
		}

		// 2️⃣ Hitung total budget
		var totalToConsume int64
		for _, r := range validReports {
			totalToConsume += r.TotalAmount
		}

		log.Println("Total to consume:", totalToConsume)

		// 3️⃣ Consume budget SEKALI
		if err := budgets.ConsumeBudgetLogic(tx, tenantID, totalToConsume); err != nil {
			return err
		}

		// 4️⃣ Update semua report
		if err := tx.Model(&models.ExpenseReport{}).
			Where("id IN ? AND tenant_id = ?", reportIDs, tenantID).
			Updates(map[string]interface{}{
				"status":         "APPROVED",
				"approved_by_id": managerID,
			}).Error; err != nil {
			return err
		}

		// 5️⃣ Insert logs
		for _, r := range validReports {
			logEntry := models.ApprovalLog{
				ExpenseReportID: &r.ID,
				UserID:          managerID,
				Action:          "APPROVE",
				Comment:         "Bulk approve by manager",
			}
			if err := tx.Create(&logEntry).Error; err != nil {
				return err
			}
		}

		log.Println("=== BULK APPROVE SUCCESS ===")
		return nil
	})
}

func BulkRejectReports(
	tenantID uuid.UUID,
	reportIDs []uuid.UUID,
	managerID uuid.UUID,
) error {

	if len(reportIDs) == 0 {
		return errors.New("no report ids provided")
	}

	// 🔥 VALIDASI UUID
	for _, id := range reportIDs {
		if id == uuid.Nil {
			return errors.New("invalid report id detected (uuid.Nil)")
		}
	}

	return configs.DB.Transaction(func(tx *gorm.DB) error {

		log.Println("=== BULK REJECT DEBUG ===")
		log.Println("TenantID:", tenantID)
		log.Println("ManagerID:", managerID)
		log.Println("ReportIDs:", reportIDs)

		// 1️⃣ Ambil semua report
		var reports []models.ExpenseReport
		if err := tx.
			Where("id IN ?", reportIDs).
			Find(&reports).Error; err != nil {
			return err
		}

		log.Println("Found Reports:", len(reports))

		if len(reports) == 0 {
			return errors.New("no reports found for given IDs")
		}

		var validReports []models.ExpenseReport
		var invalidIDs []uuid.UUID

		for _, r := range reports {

			if r.TenantID != tenantID {
				log.Println("Tenant mismatch:", r.ID)
				invalidIDs = append(invalidIDs, r.ID)
				continue
			}

			if r.Status != "SUBMITTED" {
				log.Println("Invalid status:", r.ID, r.Status)
				invalidIDs = append(invalidIDs, r.ID)
				continue
			}

			validReports = append(validReports, r)
		}

		if len(invalidIDs) > 0 {
			return fmt.Errorf("some reports invalid (tenant/status): %v", invalidIDs)
		}

		if len(validReports) != len(reportIDs) {
			return fmt.Errorf(
				"some report_ids not found (expected %d, got %d)",
				len(reportIDs),
				len(validReports),
			)
		}

		// 2️⃣ Update bulk
		if err := tx.Model(&models.ExpenseReport{}).
			Where("id IN ? AND tenant_id = ?", reportIDs, tenantID).
			Updates(map[string]interface{}{
				"status":         "REJECTED",
				"approved_by_id": managerID,
			}).Error; err != nil {
			return err
		}

		// 3️⃣ Logs
		for _, r := range validReports {
			logEntry := models.ApprovalLog{
				ExpenseReportID: &r.ID,
				UserID:          managerID,
				Action:          "REJECT",
				Comment:         "Bulk reject by manager",
			}
			if err := tx.Create(&logEntry).Error; err != nil {
				return err
			}
		}

		log.Println("=== BULK REJECT SUCCESS ===")
		return nil
	})
}
