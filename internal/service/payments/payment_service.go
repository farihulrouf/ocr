package payments

import (
	"ocr-saas-backend/internal/models"
	"ocr-saas-backend/internal/repository/payments"

	"github.com/google/uuid"
)

func GetAllpayments(
	tenantID uuid.UUID,
	page, pageSize int,
) (interface{}, error) {

	rows, total, err := payments.ListPayments(
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

func CreatePayments(
	tenantID uuid.UUID,
	name string,
) error {

	payment := &models.PaymentMethod{
		TenantID: tenantID,
		Name:     name,
	}
	return payments.CreatePayments(payment)
}

func UpdatePayments(
	tenantID, id uuid.UUID,
	name string,
) error {

	payment, err := payments.GetPaymentByID(tenantID, id)
	if err != nil {
		return err
	}

	payment.Name = name

	return payments.UpdatePayments(payment)
}

func DeleteTePayment(
	tenantID, id uuid.UUID,
) error {
	return payments.DeletePayments(tenantID, id)
}
