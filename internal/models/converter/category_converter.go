package converter

import (
	"kasir-api/internal/entity"
	"kasir-api/internal/models"
)

func CategoryToResponse(category *entity.Category) *models.CategoryResponse {
	return &models.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}
}
