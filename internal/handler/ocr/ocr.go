package ocr

import (
	"path/filepath"

	"ocr-saas-backend/internal/service/ocr"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func UploadReceipt(c *fiber.Ctx) error {
	tenantID := uuid.MustParse(c.Locals("tenant_id").(string))
	userID := uuid.MustParse(c.Locals("user_id").(string))

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "file is required"})
	}

	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".pdf" {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "invalid file type"})
	}

	savePath := "./uploads/receipts/" + uuid.New().String() + ext
	if err := c.SaveFile(file, savePath); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "failed to save file"})
	}

	receipt, err := ocr.UploadReceipt(tenantID, userID, savePath)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	// push ke queue
	if err := ocr.PushToQueue(receipt.ID); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "failed to push to OCR queue"})
	}

	// Hanya kembalikan ID
	return c.JSON(fiber.Map{
		"status": "success",
		"id":     receipt.ID,
	})
}
