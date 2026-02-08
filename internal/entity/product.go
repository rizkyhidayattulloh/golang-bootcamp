package entity

type Product struct {
	ID         int     `json:"id"`
	CategoryID int     `json:"category_id"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	Stock      int     `json:"stock"`
	Version    int     `json:"version"`

	Category *Category `json:"category"`
}
