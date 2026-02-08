package entity

type TransactionProduct struct {
	ID            int     `json:"id"`
	TransactionID int     `json:"transaction_id"`
	ProductID     int     `json:"product_id"`
	ProductName   string  `json:"product_name"`
	Quantity      int     `json:"quantity"`
	UnitPrice     float64 `json:"unit_price"`
	Subtotal      float64 `json:"subtotal"`
}
