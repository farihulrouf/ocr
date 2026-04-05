package vendors

import (
	"errors"
	"ocr-saas-backend/internal/models"
	"ocr-saas-backend/internal/repository/vendors"

	"github.com/google/uuid"
)

func GetAllVendorService(
	tenantID uuid.UUID,
	search string,
	category string,
	page, pageSize int,
) (map[string]interface{}, error) {

	rows, total, totalPages, err := vendors.ListVendor(
		tenantID, search, category, page, pageSize,
	)
	if err != nil {
		return nil, err
	}

	// Format response sesuai standar SEIDO
	return map[string]interface{}{
		"status": "success",
		"data":   rows,
		"meta": map[string]interface{}{
			"current_page": page,
			"page_size":    pageSize,
			"total_items":  total,
			"total_pages":  totalPages,
			"has_next":     page < totalPages,
			"has_prev":     page > 1,
		},
	}, nil
}

func CreateVendor(
	tenantID uuid.UUID,
	name string,
	displayName string,
	taxNumber string,
	category string,
	aliases string,
	isGlobal bool,
) error {
	// Buat objek model baru
	vendor := &models.VendorMaster{
		Base: models.Base{
			ID: uuid.New(), // Isi ID di dalam blok Base
		},
		TenantID:    &tenantID,
		Name:        name,
		DisplayName: displayName,
		TaxNumber:   taxNumber,
		Category:    category,
		Aliases:     aliases,
		IsGlobal:    isGlobal,
		IsVerified:  true, // Default verified jika diinput manual oleh admin
	}

	// Kirim ke repository
	return vendors.CreateVendor(vendor)
}

func UpdateVendorService(
	tenantID, id uuid.UUID,
	name string,
	displayName string,
	taxNumber string,
	category string,
	aliases string,
	isGlobal bool,
) error {
	// 1. Cari data vendor yang ada
	vendor, err := vendors.GetVendorByID(tenantID, id)
	if err != nil {
		return err // Vendor tidak ditemukan atau milik tenant lain
	}

	// 2. Update field (Sesuaikan dengan struct model kamu)
	vendor.Name = name
	vendor.DisplayName = displayName
	vendor.TaxNumber = taxNumber
	vendor.Category = category
	vendor.Aliases = aliases
	vendor.IsGlobal = isGlobal
	// vendor.UpdatedAt = time.Now() // Jika ada field timestamp

	// 3. Simpan perubahan ke database (Repository)
	return vendors.UpdateVendor(vendor)
}

func DeleteVendor(
	tenantID, id uuid.UUID,
) error {
	// 1. Validasi: Apakah vendor sudah punya transaksi?
	// Kita tidak ingin menghapus vendor yang sudah ada di laporan keuangan
	usage, err := vendors.CheckVendorUsage(id)
	if err != nil {
		return err
	}

	if usage > 0 {
		return errors.New("tidak bisa menghapus vendor yang sudah memiliki data transaksi")
	}

	// 2. Jika aman, panggil repository untuk hapus
	return vendors.DeleteVendor(tenantID, id)
}
