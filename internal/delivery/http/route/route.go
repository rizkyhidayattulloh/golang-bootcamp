package route

import (
	"errors"
	httpController "kasir-api/internal/delivery/http"
	"kasir-api/internal/models"
	"net/http"
	"strconv"
	"strings"
)

type RouteConfig struct {
	Server            *http.ServeMux
	ProductController *httpController.ProductController
}

func (c *RouteConfig) Setup() {
	c.Server.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		httpController.EncodeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	c.Server.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			err := c.ProductController.GetProducts(w, r)
			if err != nil {
				c.handleError(err, w)
			}
			return
		case http.MethodPost:
			err := c.ProductController.CreateProduct(w, r)
			if err != nil {
				c.handleError(err, w)
			}
		default:
			err := &models.Error{
				Status:  http.StatusMethodNotAllowed,
				Message: "Method not allowed",
			}
			c.handleError(err, w)
			return

		}
	})

	c.Server.HandleFunc("/api/products/", func(w http.ResponseWriter, r *http.Request) {
		strId := strings.TrimPrefix(r.URL.Path, "/api/products/")
		id, err := strconv.Atoi(strId)
		if err != nil {
			err := &models.Error{
				Status:  http.StatusBadRequest,
				Message: "Invalid product ID",
			}
			c.handleError(err, w)
			return
		}

		switch r.Method {
		case http.MethodGet:
			err := c.ProductController.GetProduct(w, r, id)
			if err != nil {
				c.handleError(err, w)
			}
			return
		case http.MethodPut:
			err := c.ProductController.UpdateProduct(w, r, id)
			if err != nil {
				c.handleError(err, w)
			}
			return
		case http.MethodDelete:
			err := c.ProductController.DeleteProduct(w, r, id)
			if err != nil {
				c.handleError(err, w)
			}
			return
		default:
			err := &models.Error{
				Status:  http.StatusMethodNotAllowed,
				Message: "Method not allowed",
			}
			c.handleError(err, w)
			return
		}
	})
}

func (c *RouteConfig) handleError(err error, w http.ResponseWriter) {
	var httpErr *models.Error
	if errors.As(err, &httpErr) {
		httpController.EncodeJSON(w, httpErr.Status, httpErr)
	}
}
