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
	transactionRepository := repository.NewTransactionRepository(config.DB)
	transactor := repository.NewTransactor(config.DB)

	productUseCase := usecase.NewProductUseCase(productRepository, categoryRepository, transactor, config.Validate)
	categoryUseCase := usecase.NewCategoryUseCase(categoryRepository, config.Validate)
	transactionUseCase := usecase.NewTransactionUseCase(transactionRepository, productRepository, transactor, config.Validate)
	reportUseCase := usecase.NewReportUseCase(transactionRepository)

	productController := httpController.NewProductController(productUseCase)
	categoryController := httpController.NewCategoryController(categoryUseCase)
	transactionController := httpController.NewTransactionController(transactionUseCase)
	reportController := httpController.NewReportController(reportUseCase)
	docsController := httpController.NewDocsController()

	routeConfig := route.RouteConfig{
		Server:                config.Server,
		ProductController:     productController,
		CategoryController:    categoryController,
		TransactionController: transactionController,
		ReportController:      reportController,
		DocsController:        docsController,
	}

	routeConfig.Setup()
}
