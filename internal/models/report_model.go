package models

type ReportResponse struct {
	TotalRevenue       float64          `json:"total_revenue"`
	TotalTransaction   int              `json:"total_transaction"`
	BestSellingProduct ReportTopProduct `json:"best_selling_product"`
}

type ReportTopProduct struct {
	Name      string `json:"name"`
	TotalSold int    `json:"total_sold"`
}
