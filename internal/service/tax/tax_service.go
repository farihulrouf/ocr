package tax

import (
	"ocr-saas-backend/internal/models"
	"ocr-saas-backend/internal/repository/tax"

	"github.com/google/uuid"
)

func GetTaxRates(
	tenantID uuid.UUID,
	page, pageSize int,
) (interface{}, error) {

	rows, total, err := tax.GetTaxRates(
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

func CreateTaxRate(
	tenantID uuid.UUID,
	name string,
	percentage int,
) error {

	rate := &models.TaxRate{
		TenantID:   tenantID,
		Name:       name,
		Percentage: percentage,
	}

	return tax.CreateTaxRate(rate)
}

func UpdateTaxRate(
	tenantID, id uuid.UUID,
	name string,
	percentage int,
) error {

	rate, err := tax.GetTaxRateByID(tenantID, id)
	if err != nil {
		return err
	}

	rate.Name = name
	rate.Percentage = percentage

	return tax.UpdateTaxRate(rate)
}

func DeleteTaxRate(
	tenantID, id uuid.UUID,
) error {
	return tax.DeleteTaxRate(tenantID, id)
}
