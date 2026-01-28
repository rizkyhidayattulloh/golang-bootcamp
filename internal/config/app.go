package config

import (
	"database/sql"
	httpController "kasir-api/internal/delivery/http"
	"kasir-api/internal/delivery/http/route"
	"kasir-api/internal/repository"
	"kasir-api/internal/usecase"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type BootstrapConfig struct {
	Server   *http.ServeMux
	Config   *viper.Viper
	Validate *validator.Validate
	DB       *sql.DB
}

func Bootstrap(config *BootstrapConfig) {
	productRepository := repository.NewProductRepository(config.DB)
	categoryRepository := repository.NewCategoryRepository(config.DB)

	productUseCase := usecase.NewProductUseCase(productRepository, categoryRepository, config.Validate)
	categoryUseCase := usecase.NewCategoryUseCase(categoryRepository, config.Validate)

	productController := httpController.NewProductController(productUseCase)
	categoryController := httpController.NewCategoryController(categoryUseCase)
	docsController := httpController.NewDocsController()

	routeConfig := route.RouteConfig{
		Server:             config.Server,
		ProductController:  productController,
		CategoryController: categoryController,
		DocsController:     docsController,
	}

	routeConfig.Setup()
}
