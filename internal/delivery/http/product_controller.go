package http

import "net/http"

type ProductController struct {
}

func NewProductController(server *http.ServeMux) *ProductController {
	return &ProductController{}
}

func (pc *ProductController) GetProducts(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("List of products"))
}
