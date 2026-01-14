package mapper

import (
	"ocr-saas-backend/internal/dto"
	"ocr-saas-backend/internal/models"
)

func MapReceiptToEmployeeDetailDTO(
	r *models.Receipt,
) dto.EmployeeReceiptDetailResponse {

	items := make([]dto.ReceiptDetailItem, 0, len(r.LineItems))
	for _, it := range r.LineItems {
		items = append(items, dto.ReceiptDetailItem{
			ID:          it.ID,
			Description: it.Description,
			Amount:      it.Amount,
			TaxAmount:   it.TaxAmount,
			TaxRate:     it.TaxRate,
		})
	}

	var category *dto.ReceiptDetailCategory
	if r.AccountCategory != nil {
		category = &dto.ReceiptDetailCategory{
			ID:   r.AccountCategory.ID,
			Code: r.AccountCategory.Code,
			Name: r.AccountCategory.Name,
		}
	}

	date := ""
	if r.TransactionDate != nil {
		date = r.TransactionDate.Format("2006-01-02")
	}

	return dto.EmployeeReceiptDetailResponse{
		ID:        r.ID,
		Date:      date,
		StoreName: r.StoreName,
		ImageURL:  r.ImageURL,
		Category:  category,
		Taxation:  mapTaxation(r.IsQualified),
		Amount:    r.TotalAmount,
		Status:    r.Status,
		Items:     items,
	}
}

func mapTaxation(isQualified bool) string {
	if isQualified {
		return "eligible"
	}
	return "non-eligible"
}
