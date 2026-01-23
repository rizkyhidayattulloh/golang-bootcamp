package converter

import (
	"kasir-api/internal/entity"
	"kasir-api/internal/models"
)

func ProductToResponse(product *entity.Product) *models.ProductResponse {
	return &models.ProductResponse{
		ID:    product.ID,
		Name:  product.Name,
		Price: product.Price,
		Stock: product.Stock,
	}
}
