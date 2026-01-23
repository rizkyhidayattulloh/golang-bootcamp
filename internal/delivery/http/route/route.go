package route

import (
	httpController "kasir-api/internal/delivery/http"
	"net/http"
)

type RouteConfig struct {
	Server            *http.ServeMux
	ProductController *httpController.ProductController
}

func (c *RouteConfig) Setup() {
	c.Server.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	c.Server.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			c.ProductController.GetProducts(w, r)
			return
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return

		}
	})
}
