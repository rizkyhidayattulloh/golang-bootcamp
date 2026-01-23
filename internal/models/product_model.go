package models

type ProductResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

type CreateProductRequest struct {
	Name  string `json:"name" validate:"required"`
	Price int    `json:"price" validate:"required,min=0"`
	Stock int    `json:"stock" validate:"required,min=0"`
}

type UpdateProductRequest struct {
	ID    int    `json:"id" validate:"required"`
	Name  string `json:"name" validate:"required"`
	Price int    `json:"price" validate:"required,min=0"`
	Stock int    `json:"stock" validate:"required,min=0"`
}
