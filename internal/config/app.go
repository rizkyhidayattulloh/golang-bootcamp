package config

import (
	httpController "kasir-api/internal/delivery/http"
	"kasir-api/internal/delivery/http/route"
	"kasir-api/internal/repository"
	"kasir-api/internal/usecase"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type BootstrapConfig struct {
	Server   *http.ServeMux
	Config   *Config
	Validate *validator.Validate
}

func Bootstrap(config *BootstrapConfig) {
	productRepository := repository.NewProductRepository()

	productUseCase := usecase.NewProductUseCase(productRepository, config.Validate)

	productController := httpController.NewProductController(productUseCase)

	routeConfig := route.RouteConfig{
		Server:            config.Server,
		ProductController: productController,
	}

	routeConfig.Setup()
}
