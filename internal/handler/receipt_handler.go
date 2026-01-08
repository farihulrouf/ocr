package handler

import (
	"ocr-saas-backend/internal/service"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetMyReceipts(c *fiber.Ctx) error {
	tenantID, ok := c.Locals("tenant_id").(string)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "unauthorized"})
	}

	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "unauthorized"})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	if pageSize < 1 {
		pageSize = 10
	}

	q := c.Query("q")
	status := c.Query("status")
	sort := c.Query("sort")

	response, err := service.GetMyReceipts(
		uuid.MustParse(tenantID),
		uuid.MustParse(userID),
		page,
		pageSize,
		q,
		status,
		sort,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.JSON(response)
}

func GetAllReceipts(c *fiber.Ctx) error {
	tenantID, ok := c.Locals("tenant_id").(string)
	if !ok {
		return c.Status(401).JSON(fiber.Map{
			"status":  "error",
			"message": "unauthorized",
		})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	if pageSize < 1 {
		pageSize = 10
	}

	q := c.Query("q")
	status := c.Query("status")
	sort := c.Query("sort")

	response, err := service.GetAllReceipts(
		uuid.MustParse(tenantID),
		page,
		pageSize,
		q,
		status,
		sort,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.JSON(response)
}

func GetReceiptDetail(c *fiber.Ctx) error {
	tenantIDStr, ok := c.Locals("tenant_id").(string)
	if !ok {
		return fiber.NewError(
			fiber.StatusUnauthorized,
			"invalid tenant context",
		)
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return fiber.NewError(
			fiber.StatusUnauthorized,
			"invalid tenant id",
		)
	}

	receiptID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"invalid receipt id",
		)
	}

	result, err := service.GetReceiptDetail(tenantID, receiptID)
	if err != nil {
		return fiber.NewError(
			fiber.StatusNotFound,
			err.Error(),
		)
	}

	return c.JSON(result)
}

func ConfirmReceipt(c *fiber.Ctx) error {
	tenantIDStr, ok := c.Locals("tenant_id").(string)
	if !ok {
		return fiber.NewError(
			fiber.StatusUnauthorized,
			"invalid tenant context",
		)
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return fiber.NewError(
			fiber.StatusUnauthorized,
			"invalid tenant id",
		)
	}

	receiptID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"invalid receipt id",
		)
	}

	// request body
	var req struct {
		Total int64  `json:"total"`
		Date  string `json:"date"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"invalid request body",
		)
	}

	if req.Total <= 0 {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"total must be greater than zero",
		)
	}

	transactionDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"invalid date format (use YYYY-MM-DD)",
		)
	}
	// call service
	err = service.ConfirmReceipt(
		tenantID,
		receiptID,
		req.Total,
		transactionDate,
	)
	if err != nil {
		switch err {
		case service.ErrReceiptNotFound:
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		case service.ErrReceiptAlreadyFinal:
			return fiber.NewError(fiber.StatusConflict, err.Error())
		case service.ErrInvalidTotalAmount:
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(fiber.Map{
		"status": "confirmed",
	})
}

func DeleteReceipt(c *fiber.Ctx) error {
	tenantIDStr, ok := c.Locals("tenant_id").(string)
	if !ok {
		return fiber.NewError(
			fiber.StatusUnauthorized,
			"invalid tenant context",
		)
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return fiber.NewError(
			fiber.StatusUnauthorized,
			"invalid tenant id",
		)
	}

	receiptID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"invalid receipt id",
		)
	}

	err = service.DeleteReceiptManager(tenantID, receiptID)
	if err != nil {
		if err == service.ErrReceiptNotFound {
			return fiber.NewError(
				fiber.StatusNotFound,
				err.Error(),
			)
		}
		return fiber.NewError(
			fiber.StatusInternalServerError,
			err.Error(),
		)
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "receipt deleted",
	})
}

func BulkDeleteReceipts(c *fiber.Ctx) error {
	tenantIDStr, ok := c.Locals("tenant_id").(string)
	if !ok {
		return fiber.NewError(
			fiber.StatusUnauthorized,
			"invalid tenant context",
		)
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return fiber.NewError(
			fiber.StatusUnauthorized,
			"invalid tenant id",
		)
	}

	// request body
	var req struct {
		IDs []string `json:"ids"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"invalid request body",
		)
	}

	if len(req.IDs) == 0 {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"ids cannot be empty",
		)
	}

	ids := make([]uuid.UUID, 0, len(req.IDs))
	for _, idStr := range req.IDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return fiber.NewError(
				fiber.StatusBadRequest,
				"invalid receipt id",
			)
		}
		ids = append(ids, id)
	}

	deleted, err := service.BulkDeleteReceiptsManager(
		tenantID,
		ids,
	)
	if err != nil {
		switch err {
		case service.ErrNoReceiptDeleted:
			return fiber.NewError(
				fiber.StatusNotFound,
				err.Error(),
			)
		default:
			return fiber.NewError(
				fiber.StatusInternalServerError,
				err.Error(),
			)
		}
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"deleted": deleted,
	})
}
