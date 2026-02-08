package entity

import "time"

type Transaction struct {
	ID          int       `json:"id"`
	TotalAmount float64   `json:"total_amount"`
	CreatedAt   time.Time `json:"created_at"`

	Products []TransactionProduct `json:"products"`
}
