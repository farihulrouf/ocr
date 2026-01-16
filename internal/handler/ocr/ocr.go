package ocr

import (
	"ocr-saas-backend/internal/service/ocr"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func UploadReceipt(c *fiber.Ctx) error {
	// ===== ambil context =====
	tenantID := uuid.MustParse(c.Locals("tenant_id").(string))
	userID := uuid.MustParse(c.Locals("user_id").(string))

	// ===== ambil file =====
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "file is required",
		})
	}

	// ===== validasi extension =====
	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".pdf" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid file type",
		})
	}

	// ===== simpan file (sementara local dulu) =====
	savePath := "./uploads/receipts/" + uuid.New().String() + ext

	if err := c.SaveFile(file, savePath); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to save file",
		})
	}

	// ===== PANGGIL SERVICE (INTI) =====
	receipt, err := ocr.UploadReceipt(
		tenantID,
		userID,
		savePath,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	go ocr.ProcessOCR(receipt.ID)
	// ===== response =====
	return c.JSON(fiber.Map{
		"status": "success",
		"data":   receipt,
	})
}
