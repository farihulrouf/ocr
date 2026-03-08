package ocr

import (
	"path/filepath"

	"ocr-saas-backend/configs"
	"ocr-saas-backend/internal/service/ocr"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func UploadReceipt(c *fiber.Ctx) error {
	tenantID := uuid.MustParse(c.Locals("tenant_id").(string))
	userID := uuid.MustParse(c.Locals("user_id").(string))

	// ambil file dari form
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "file is required"})
	}

	// 🔹 BATASI SIZE FILE (max 3MB)
	const MaxFileSize = 3 * 1024 * 1024 // 3MB
	if file.Size > MaxFileSize {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "file too large, max 3MB",
		})
	}

	// validasi ekstensi
	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "invalid file type"})
	}

	// buka file untuk upload
	f, err := file.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "cannot open file"})
	}
	defer f.Close()

	// generate UUID untuk object key
	objectKey := "receipts/" + uuid.New().String() + ext

	// inisialisasi MinIO client
	//storageClient := storage.NewMinioStorage(configs.MinioConfig)

	// upload ke MinIO
	//if err := storageClient.Upload(c.Context(), objectKey, f, file.Size, file.Header.Get("Content-Type")); err != nil {
	//	return c.Status(500).JSON(fiber.Map{"status": "error", "message": "failed to upload to MinIO"})
	//}

	// upload ke S3
	err = configs.S3Client.Upload(
		c.Context(),
		objectKey,
		f,
		file.Size,
		file.Header.Get("Content-Type"),
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to upload to S3",
		})
	}

	// simpan objectKey di DB (fungsi ocr.UploadReceipt harus menerima objectKey)
	receipt, err := ocr.UploadReceipt(tenantID, userID, objectKey)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	// push ke queue untuk OCR
	if err := ocr.PushToQueue(receipt.ID); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "failed to push to OCR queue"})
	}

	// kembalikan ID saja
	return c.JSON(fiber.Map{
		"status": "success",
		"id":     receipt.ID,
	})
}
