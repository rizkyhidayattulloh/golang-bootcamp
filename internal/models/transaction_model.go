package models

type CheckoutRequest struct {
	Products []CheckoutProductRequest `json:"products" validate:"required,min=1,dive"`
}

type CheckoutProductRequest struct {
	ProductID int `json:"product_id" validate:"required"`
	Quantity  int `json:"quantity" validate:"required,min=1"`
}

type CheckoutResponse struct {
	TransactionID int                     `json:"transaction_id"`
	TotalAmount   float64                 `json:"total_amount"`
	Products      []CheckoutProductDetail `json:"products"`
}

type CheckoutProductDetail struct {
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}
