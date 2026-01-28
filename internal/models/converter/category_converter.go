package converter

import (
	"kasir-api/internal/entity"
	"kasir-api/internal/models"
)

func CategoryToResponse(category *entity.Category) *models.CategoryResponse {
	var desc *string
	if category.Description.Valid {
		desc = &category.Description.String
	}

	return &models.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: desc,
	}
}
