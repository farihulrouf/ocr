package dto

import "github.com/google/uuid"

type AccountCategoryRequest struct {
	Code string `json:"code" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type AccountCategoryResponse struct {
	ID   uuid.UUID `json:"id"`
	Code string    `json:"code"`
	Name string    `json:"name"`
}
