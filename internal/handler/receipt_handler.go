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
	// tenant
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

	// user
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok {
		return fiber.NewError(
			fiber.StatusUnauthorized,
			"invalid user context",
		)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return fiber.NewError(
			fiber.StatusUnauthorized,
			"invalid user id",
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
		userID,
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

func BulkRestoreReceipts(c *fiber.Ctx) error {
	// ===== TENANT =====
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

	// ===== USER =====
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok {
		return fiber.NewError(
			fiber.StatusUnauthorized,
			"invalid user context",
		)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return fiber.NewError(
			fiber.StatusUnauthorized,
			"invalid user id",
		)
	}

	// ===== REQUEST BODY =====
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

	// ===== SERVICE =====
	restored, err := service.BulkRestoreReceiptsManager(
		tenantID,
		userID,
		ids,
	)
	if err != nil {
		if err == service.ErrNoReceiptRestored {
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
		"status":   "success",
		"restored": restored,
	})
}

// internal/handler/receipt_handler.go

func BulkApproveReceipts(c *fiber.Ctx) error {
	return bulkApproveReject(c, "APPROVE")
}

func BulkRejectReceipts(c *fiber.Ctx) error {
	return bulkApproveReject(c, "REJECT")
}

func bulkApproveReject(c *fiber.Ctx, action string) error {
	// tenant
	tenantIDStr, ok := c.Locals("tenant_id").(string)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid tenant context")
	}
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid tenant id")
	}

	// user
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid user context")
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid user id")
	}

	// request
	var req struct {
		IDs []string `json:"ids"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if len(req.IDs) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "ids cannot be empty")
	}

	ids := make([]uuid.UUID, 0, len(req.IDs))
	for _, idStr := range req.IDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid receipt id")
		}
		ids = append(ids, id)
	}

	// mapping action
	var auditAction string
	switch action {
	case "APPROVE":
		auditAction = "BULK_APPROVE_RECEIPT"
	case "REJECT":
		auditAction = "BULK_REJECT_RECEIPT"
	default:
		return fiber.NewError(fiber.StatusBadRequest, "invalid action")
	}

	updated, err := service.BulkApproveRejectReceipts(
		tenantID,
		userID,
		ids,
		action,
		auditAction,
	)
	if err != nil {
		if err == service.ErrNoReceiptUpdated {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"updated": updated,
	})
}

func BulkUpdateReceiptCategory(c *fiber.Ctx) error {

	tenantID := uuid.MustParse(c.Locals("tenant_id").(string))
	userID := uuid.MustParse(c.Locals("user_id").(string))
	role := c.Locals("role").(string)

	var req struct {
		IDs   []string `json:"ids"`
		CatID string   `json:"cat_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}

	categoryID, err := uuid.Parse(req.CatID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid category id")
	}

	var ids []uuid.UUID
	for _, s := range req.IDs {
		id, err := uuid.Parse(s)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid receipt id")
		}
		ids = append(ids, id)
	}

	updated, err := service.BulkUpdateReceiptCategory(
		tenantID,
		userID,
		role,
		ids,
		categoryID,
	)

	if err != nil {
		switch err {
		case service.ErrForbidden:
			return fiber.NewError(fiber.StatusForbidden, err.Error())
		case service.ErrCategoryNotBelongTenant:
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		case service.ErrNoReceiptUpdated:
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"updated": updated,
	})
}

func AddReceiptItem(c *fiber.Ctx) error {
	// tenant
	tenantID := uuid.MustParse(c.Locals("tenant_id").(string))

	// receipt id
	receiptID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid receipt id")
	}

	// body
	var req struct {
		Name  string `json:"name"`
		Price int64  `json:"price"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}

	if req.Name == "" || req.Price <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "name and price required")
	}

	itemID, err := service.AddReceiptItem(
		c.Context(),
		tenantID,
		receiptID,
		req.Name,
		req.Price,
	)
	if err != nil {
		if err == service.ErrReceiptNotFound {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"item_id": itemID,
	})
}

func UpdateReceiptItem(c *fiber.Ctx) error {
	itemIDParam := c.Params("itemId")
	itemID, err := strconv.ParseUint(itemIDParam, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid item id")
	}

	var req struct {
		Price int64 `json:"price"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if req.Price <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "price must be greater than zero")
	}

	// Ambil user_id dari context Fiber
	userIDStr := c.Locals("user_id")
	if userIDStr == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "user_id missing in context")
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid user_id")
	}

	// Kirim user_id ke service
	if err := service.UpdateReceiptItem(
		c.Context(),
		uint(itemID),
		req.Price,
		userID,
	); err != nil {

		switch err {
		case service.ErrItemNotFound:
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		case service.ErrReceiptNotEditable:
			return fiber.NewError(fiber.StatusConflict, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(fiber.Map{"message": "Item updated"})
}

func DeleteReceiptItem(c *fiber.Ctx) error {

	itemID, err := strconv.ParseUint(c.Params("itemId"), 10, 64)
	if err != nil {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"invalid item id",
		)
	}

	err = service.DeleteReceiptItem(
		c.Context(),
		uint(itemID),
	)

	if err != nil {
		switch err {
		case service.ErrItemNotFound:
			return fiber.NewError(fiber.StatusNotFound, err.Error())

		case service.ErrReceiptNotEditable:
			return fiber.NewError(fiber.StatusConflict, err.Error())

		default:
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(fiber.Map{
		"message": "Item deleted",
	})
}

func GetMyReceiptDetail(c *fiber.Ctx) error {

	tenantID := uuid.MustParse(c.Locals("tenant_id").(string))
	userID := uuid.MustParse(c.Locals("user_id").(string))
	receiptID := uuid.MustParse(c.Params("id"))

	result, err := service.GetMyReceiptDetail(
		tenantID,
		userID,
		receiptID,
	)
	if err != nil {
		switch err {
		case service.ErrForbidden:
			return fiber.NewError(fiber.StatusForbidden, err.Error())
		default:
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
	}

	return c.JSON(result)
}

func UpdateReceipt(c *fiber.Ctx) error {
	tenantID := uuid.MustParse(c.Locals("tenant_id").(string))
	receiptID := uuid.MustParse(c.Params("id"))

	var req struct {
		StoreName string `json:"store_name"`
		Date      string `json:"date"`
		Total     *int64 `json:"total"` // optional
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(400, "invalid body")
	}

	var txDate *time.Time
	if req.Date != "" {
		d, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			return fiber.NewError(400, "invalid date")
		}
		txDate = &d
	}

	err := service.UpdateReceipt(
		tenantID,
		receiptID,
		req.StoreName,
		txDate,
		req.Total,
	)
	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	return c.JSON(fiber.Map{"status": "updated"})
}
