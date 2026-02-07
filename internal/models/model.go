package models

import "strconv"

type WebResponse[T any] struct {
	Data T                   `json:"data"`
	Meta *PaginationResponse `json:"meta,omitempty"`
}

type PaginationResponse struct {
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
	TotalItems  int `json:"total_items"`
	TotalPages  int `json:"total_pages"`
}

type PaginationRequest struct {
	Page  int `json:"page" validate:"min=1"`
	Limit int `json:"limit" validate:"min=1,max=100"`
}

func NewPaginationRequest(page string, limit string) PaginationRequest {
	p, err := strconv.Atoi(page)
	if err != nil || p < 1 {
		p = 1
	}

	l, err := strconv.Atoi(limit)
	if err != nil || l < 1 {
		l = 10
	} else if l > 100 {
		l = 100
	}

	return PaginationRequest{
		Page:  p,
		Limit: l,
	}
}

func NewPaginationResponse(currentPage int, perPage int, totalItems int) *PaginationResponse {
	totalPages := (totalItems + perPage - 1) / perPage

	return &PaginationResponse{
		CurrentPage: currentPage,
		PerPage:     perPage,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
	}
}
