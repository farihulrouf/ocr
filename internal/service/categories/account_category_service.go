package categories

import (
	"errors"
	"ocr-saas-backend/internal/dto"
	"ocr-saas-backend/internal/models"
	"ocr-saas-backend/internal/repository/categories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func ListCategories(
	tenantID uuid.UUID,
	page, pageSize int,
) (map[string]interface{}, error) {

	rows, total, err := categories.ListCategories(tenantID, page, pageSize)
	if err != nil {
		return nil, err
	}

	resp := make([]dto.AccountCategoryResponse, 0, len(rows))
	for _, r := range rows {
		resp = append(resp, dto.AccountCategoryResponse{
			ID:   r.ID,
			Code: r.Code,
			Name: r.Name,
		})
	}

	return map[string]interface{}{
		"status": "success",
		"data":   resp,
		"meta": map[string]interface{}{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	}, nil
}

func CreateCategory(
	tenantID uuid.UUID,
	req dto.AccountCategoryRequest,
) error {

	cat := models.AccountCategory{
		TenantID: tenantID,
		Code:     req.Code,
		Name:     req.Name,
	}

	return categories.CreateCategory(&cat)
}

func UpdateCategory(
	tenantID, id uuid.UUID,
	req dto.AccountCategoryRequest,
) error {

	cat, err := categories.GetCategoryByID(tenantID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("category not found")
		}
		return err
	}

	cat.Code = req.Code
	cat.Name = req.Name

	return categories.UpdateCategory(cat)
}

func DeleteCategory(
	tenantID, id uuid.UUID,
) error {

	cat, err := categories.GetCategoryByID(tenantID, id)
	if err != nil {
		return err
	}

	return categories.DeleteCategory(cat)
}
