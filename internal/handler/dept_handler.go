package handler

import (
	"ocr-saas-backend/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func ListDepartments(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	q := c.Query("q", "")
	sort := c.Query("sort", "")

	result, err := service.GetAllDepartments(page, pageSize, q, sort)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}

	return c.Status(200).JSON(result)
}

type CreateDeptRequest struct {
	Name string `json:"name"`
}

func CreateDepartment(c *fiber.Ctx) error {
	var body CreateDeptRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if body.Name == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Name is required",
		})
	}

	tenantID, _ := uuid.Parse(c.Locals("tenant_id").(string))

	dept, err := service.CreateDepartment(tenantID, body.Name)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create department",
		})
	}

	return c.JSON(fiber.Map{
		"id":   dept.ID,
		"name": dept.Name,
	})
}

func GetDepartmentDetailHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")

	deptID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid department ID",
		})
	}

	result, err := service.GetDepartmentDetail(deptID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Department not found",
		})
	}

	return c.JSON(result)
}

func UpdateDepartment(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing department ID",
		})
	}

	var req struct {
		Name string `json:"name"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid JSON",
		})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Name is required",
		})
	}

	// call service
	if err := service.UpdateDepartment(id, req.Name); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Updated",
	})
}

func DeleteDepartment(c *fiber.Ctx) error {
	idParam := c.Params("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid department id",
		})
	}

	err = service.DeleteDepartment(id)
	if err != nil {
		if err.Error() == "department not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if err.Error() == "department still has users" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "cannot delete department with active users",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "department deleted successfully",
	})
}
