package service

import "ocr-saas-backend/internal/repository"

func GetAllDepartments(page, pageSize int, q, sort string) (interface{}, error) {
	data, total, err := repository.GetAllDepartments(page, pageSize, q, sort)
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"data": data,
		"meta": map[string]interface{}{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
		"status": "success",
	}

	return response, nil
}
