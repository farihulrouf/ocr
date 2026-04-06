package vendors

import (
	"errors"
	"math"
	"ocr-saas-backend/configs"
	"ocr-saas-backend/internal/models"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func ListVendor(
	tenantID uuid.UUID,
	search string,
	category string,
	page, pageSize int,
) ([]models.VendorMaster, int64, int, error) {

	var rows []models.VendorMaster
	var total int64

	// 1. Inisialisasi Query dengan akses Tenant + Global Vendor
	db := configs.DB.Model(&models.VendorMaster{}).
		Where("(tenant_id = ? OR is_global = ?)", tenantID, true)

	// 2. Filter Search (Case-Insensitive untuk Postgres/SQLite)
	if search != "" {
		s := "%" + strings.ToLower(search) + "%"
		// Mencari di Nama Legal, Nama Display, NPWP, dan Aliases OCR
		db = db.Where("(LOWER(name) LIKE ? OR LOWER(display_name) LIKE ? OR tax_number LIKE ? OR LOWER(aliases) LIKE ?)",
			s, s, s, s)
	}

	// 3. Filter Category
	if category != "" && category != "All" {
		db = db.Where("category = ?", category)
	}

	// 4. Hitung Total Data SETELAH filter (untuk pagination)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	// 5. Kalkulasi Total Halaman
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	if totalPages == 0 {
		totalPages = 1
	}

	// 6. Eksekusi Query dengan Limit & Offset
	offset := (page - 1) * pageSize
	err := db.
		Order("is_global DESC, display_name ASC"). // Global di atas, lalu urut abjad
		Limit(pageSize).
		Offset(offset).
		Find(&rows).Error

	return rows, total, totalPages, err
}

// Di Repository/Service Backend
func GetVendorByID(tenantID uuid.UUID, id uuid.UUID) (*models.VendorMaster, error) {
	var vendor models.VendorMaster
	// Tambahkan kondisi OR is_global = true agar tidak "Not Found" saat edit vendor global
	err := configs.DB.Where("id = ? AND (tenant_id = ? OR is_global = true)", id, tenantID).First(&vendor).Error
	return &vendor, err
}

func CreateVendor(cat *models.VendorMaster) error {
	return configs.DB.Create(cat).Error
}

func UpdateVendor(cat *models.VendorMaster) error {
	return configs.DB.Save(cat).Error
}

// DeleteVendorRepo mengeksekusi penghapusan di level DB
func DeleteVendor(tenantID, id uuid.UUID) error {
	// 1. Cek apakah vendor ada dan milik tenant (Security Check)
	var vendor models.VendorMaster
	err := configs.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&vendor).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("vendor tidak ditemukan atau anda tidak memiliki akses")
		}
		return err
	}

	// 2. Jalankan Soft Delete (karena ada gorm.DeletedAt di Base struct)
	return configs.DB.Delete(&vendor).Error
}

// CheckVendorUsage cek apakah vendor sudah dipakai di transaksi lain
func CheckVendorUsage(vendorID uuid.UUID) (int64, error) {
	var count int64

	// Ganti models.Expense menjadi models.Receipt karena itu yang ada di models kamu
	err := configs.DB.Model(&models.Receipt{}).Where("vendor_id = ?", vendorID).Count(&count).Error

	return count, err
}
