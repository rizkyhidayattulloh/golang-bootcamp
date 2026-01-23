package usecase

import (
	"kasir-api/internal/entity"
	"kasir-api/internal/models"
	"kasir-api/internal/models/converter"
	"kasir-api/internal/repository"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type ProductUseCase struct {
	ProductRepository *repository.ProductRepository
	Validate          *validator.Validate
}

func NewProductUseCase(productRepository *repository.ProductRepository, validate *validator.Validate) *ProductUseCase {
	return &ProductUseCase{
		ProductRepository: productRepository,
		Validate:          validate,
	}
}

func (p *ProductUseCase) List() ([]models.ProductResponse, error) {
	products, err := p.ProductRepository.GetAll()
	if err != nil {
		return nil, &models.Error{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	response := make([]models.ProductResponse, len(products))
	for i, product := range products {
		response[i] = *converter.ProductToResponse(&product)
	}

	return response, nil
}

func (p *ProductUseCase) Get(id int) (*models.ProductResponse, error) {
	product, err := p.ProductRepository.FindByID(id)
	if err != nil {
		return nil, &models.Error{
			Status:  http.StatusNotFound,
			Message: err.Error(),
		}
	}

	return converter.ProductToResponse(product), nil
}

func (p *ProductUseCase) Create(request *models.CreateProductRequest) (*models.ProductResponse, error) {
	err := p.Validate.Struct(request)
	if err != nil {
		return nil, &models.Error{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	product := &entity.Product{
		Name:  request.Name,
		Price: request.Price,
		Stock: request.Stock,
	}

	createdProduct, err := p.ProductRepository.Create(product)
	if err != nil {
		return nil, &models.Error{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	response := converter.ProductToResponse(createdProduct)
	return response, nil
}

func (p *ProductUseCase) Update(request *models.UpdateProductRequest) (*models.ProductResponse, error) {

	err := p.Validate.Struct(request)
	if err != nil {
		return nil, &models.Error{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	_, err = p.ProductRepository.FindByID(request.ID)
	if err != nil {
		return nil, &models.Error{
			Status:  http.StatusNotFound,
			Message: err.Error(),
		}
	}

	product := &entity.Product{
		ID:    request.ID,
		Name:  request.Name,
		Price: request.Price,
		Stock: request.Stock,
	}

	updatedProduct, err := p.ProductRepository.Update(product)
	if err != nil {
		return nil, &models.Error{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	response := converter.ProductToResponse(updatedProduct)
	return response, nil
}

func (p *ProductUseCase) Delete(id int) error {
	_, err := p.ProductRepository.FindByID(id)
	if err != nil {
		return &models.Error{
			Status:  http.StatusNotFound,
			Message: err.Error(),
		}
	}

	err = p.ProductRepository.Delete(id)
	if err != nil {
		return &models.Error{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil
}
