package disbursement

import (
	"ocr-saas-backend/internal/service/disbursement"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func ExecutePayment(c *fiber.Ctx) error {
	// Ambil ID dari Locals (Format string ke UUID sesuai gaya Mas)
	tenantID := uuid.MustParse(c.Locals("tenant_id").(string))
	financeID := uuid.MustParse(c.Locals("user_id").(string))

	var req struct {
		ReportID        string `json:"report_id"`
		Amount          int64  `json:"amount"`
		ReferenceNumber string `json:"reference_number"`
		ProofImageURL   string `json:"proof_image_url"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(400, "invalid body")
	}

	reportID, err := uuid.Parse(req.ReportID)
	if err != nil {
		return fiber.NewError(400, "invalid report id")
	}

	// Panggil Service
	err = disbursement.ProcessPayment(
		reportID,
		financeID,
		tenantID,
		req.Amount,
		req.ReferenceNumber,
		req.ProofImageURL,
	)

	if err != nil {
		return fiber.NewError(422, err.Error())
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Payment executed",
	})
}
