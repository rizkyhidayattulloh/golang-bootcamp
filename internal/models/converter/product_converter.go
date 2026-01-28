package converter

import (
	"kasir-api/internal/entity"
	"kasir-api/internal/models"
)

func ProductToResponse(product *entity.Product) *models.ProductResponse {
	response := &models.ProductResponse{
		ID:    product.ID,
		Name:  product.Name,
		Price: product.Price,
		Stock: product.Stock,
	}

	if product.Category != nil {
		response.Category = &models.ProductCategoryResponse{
			ID:   product.Category.ID,
			Name: product.Category.Name,
		}
	}

	return response
}
