package vendors

import (
	"ocr-saas-backend/internal/models"
	"ocr-saas-backend/internal/repository/vendors"

	"github.com/google/uuid"
)

func GetAllVendor(
	tenantID uuid.UUID,
	page, pageSize int,
) (interface{}, error) {

	rows, total, err := vendors.ListVendor(
		tenantID, page, pageSize,
	)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"data": rows,
		"meta": map[string]interface{}{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
		"status": "success",
	}, nil
}

func CreateVendor(
	tenantID uuid.UUID,
	name string,
	TaxNumber string,
) error {

	vendor := &models.VendorMaster{
		TenantID:  tenantID,
		Name:      name,
		TaxNumber: TaxNumber,
	}

	return vendors.CreateVendor(vendor)
}

func UpdateVendor(
	tenantID, id uuid.UUID,
	name string,
	TaxNumber string,
) error {

	vendor, err := vendors.GetVendorByID(tenantID, id)
	if err != nil {
		return err
	}

	vendor.Name = name
	vendor.TaxNumber = TaxNumber

	return vendors.UpdateVendor(vendor)
}

func DeleteVendor(
	tenantID, id uuid.UUID,
) error {
	return vendors.DeleteVendor(tenantID, id)
}
