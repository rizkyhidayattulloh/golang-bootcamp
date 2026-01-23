package models

type CategoryResponse struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type CreateCategoryRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description" validate:"omitempty"`
}

type UpdateCategoryRequest struct {
	ID          int     `json:"id" validate:"required"`
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description" validate:"omitempty"`
}
