package models

type ProductResponse struct {
	ID       int                      `json:"id"`
	Name     string                   `json:"name"`
	Price    float64                  `json:"price"`
	Stock    int                      `json:"stock"`
	Category *ProductCategoryResponse `json:"category,omitempty"`
}

type ProductCategoryResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CreateProductRequest struct {
	CategoryID int     `json:"category_id" validate:"required"`
	Name       string  `json:"name" validate:"required"`
	Price      float64 `json:"price" validate:"required,min=0"`
	Stock      int     `json:"stock" validate:"required,min=0"`
}

type UpdateProductRequest struct {
	CategoryID int     `json:"category_id" validate:"required"`
	ID         int     `json:"id" validate:"required"`
	Name       string  `json:"name" validate:"required"`
	Price      float64 `json:"price" validate:"required,min=0"`
	Stock      int     `json:"stock" validate:"required,min=0"`
}
