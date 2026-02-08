package usecase

import (
	"context"
	"kasir-api/internal/entity"
	"kasir-api/internal/models"
	"kasir-api/internal/models/converter"
	"kasir-api/internal/repository"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type ProductUseCase struct {
	ProductRepository  *repository.ProductRepository
	CategoryRepository *repository.CategoryRepository
	Transactor         *repository.Transactor
	Validate           *validator.Validate
}

func NewProductUseCase(productRepository *repository.ProductRepository, categoryRepository *repository.CategoryRepository, transactor *repository.Transactor, validate *validator.Validate) *ProductUseCase {
	return &ProductUseCase{
		ProductRepository:  productRepository,
		CategoryRepository: categoryRepository,
		Transactor:         transactor,
		Validate:           validate,
	}
}

func (p *ProductUseCase) List(ctx context.Context, request *models.SearchProductRequest) ([]models.ProductResponse, *models.PaginationResponse, error) {
	products, total, err := p.ProductRepository.GetAndCountAll(ctx, request)
	if err != nil {
		return nil, nil, &models.Error{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	response := make([]models.ProductResponse, len(products))
	for i, product := range products {
		response[i] = *converter.ProductToResponse(&product)
	}

	pagination := models.NewPaginationResponse(request.Page, request.Limit, total)

	return response, pagination, nil
}

func (p *ProductUseCase) Get(ctx context.Context, id int) (*models.ProductResponse, error) {
	product, err := p.ProductRepository.FindByID(ctx, id)
	if err != nil {
		return nil, &models.Error{
			Status:  http.StatusNotFound,
			Message: err.Error(),
		}
	}

	return converter.ProductToResponse(product), nil
}

func (p *ProductUseCase) Create(ctx context.Context, request *models.CreateProductRequest) (*models.ProductResponse, error) {
	err := p.Validate.Struct(request)
	if err != nil {
		return nil, &models.Error{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	product := &entity.Product{
		Name:       request.Name,
		Price:      request.Price,
		Stock:      request.Stock,
		CategoryID: request.CategoryID,
	}

	var response *models.ProductResponse
	err = p.Transactor.WithTx(ctx, func(txCtx context.Context) error {
		_, err := p.CategoryRepository.FindByID(txCtx, request.CategoryID)
		if err != nil {
			return &models.Error{
				Status:  http.StatusBadRequest,
				Message: "invalid category id",
			}
		}

		createdProduct, err := p.ProductRepository.Create(txCtx, product)
		if err != nil {
			return &models.Error{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			}
		}

		response = converter.ProductToResponse(createdProduct)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (p *ProductUseCase) Update(ctx context.Context, request *models.UpdateProductRequest) (*models.ProductResponse, error) {

	err := p.Validate.Struct(request)
	if err != nil {
		return nil, &models.Error{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	product := &entity.Product{
		ID:         request.ID,
		Name:       request.Name,
		Price:      request.Price,
		Stock:      request.Stock,
		CategoryID: request.CategoryID,
	}

	var response *models.ProductResponse
	err = p.Transactor.WithTx(ctx, func(txCtx context.Context) error {
		_, err := p.ProductRepository.FindByID(txCtx, request.ID)
		if err != nil {
			return &models.Error{
				Status:  http.StatusNotFound,
				Message: err.Error(),
			}
		}

		_, err = p.CategoryRepository.FindByID(txCtx, request.CategoryID)
		if err != nil {
			return &models.Error{
				Status:  http.StatusBadRequest,
				Message: "invalid category id",
			}
		}

		updatedProduct, err := p.ProductRepository.Update(txCtx, product)
		if err != nil {
			return &models.Error{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			}
		}

		response = converter.ProductToResponse(updatedProduct)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (p *ProductUseCase) Delete(ctx context.Context, id int) error {
	_, err := p.ProductRepository.FindByID(ctx, id)
	if err != nil {
		return &models.Error{
			Status:  http.StatusNotFound,
			Message: err.Error(),
		}
	}

	err = p.ProductRepository.Delete(ctx, id)
	if err != nil {
		return &models.Error{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil
}
