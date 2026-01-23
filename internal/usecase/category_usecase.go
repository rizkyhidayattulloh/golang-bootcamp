package usecase

import (
	"kasir-api/internal/entity"
	"kasir-api/internal/models"
	"kasir-api/internal/models/converter"
	"kasir-api/internal/repository"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type CategoryUseCase struct {
	CategoryRepository *repository.CategoryRepository
	Validate           *validator.Validate
}

func NewCategoryUseCase(CategoryRepository *repository.CategoryRepository, validate *validator.Validate) *CategoryUseCase {
	return &CategoryUseCase{
		CategoryRepository: CategoryRepository,
		Validate:           validate,
	}
}

func (p *CategoryUseCase) List() ([]models.CategoryResponse, error) {
	categories, err := p.CategoryRepository.GetAll()
	if err != nil {
		return nil, &models.Error{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	response := make([]models.CategoryResponse, len(categories))
	for i, Category := range categories {
		response[i] = *converter.CategoryToResponse(&Category)
	}

	return response, nil
}

func (p *CategoryUseCase) Get(id int) (*models.CategoryResponse, error) {
	Category, err := p.CategoryRepository.FindByID(id)
	if err != nil {
		return nil, &models.Error{
			Status:  http.StatusNotFound,
			Message: err.Error(),
		}
	}

	return converter.CategoryToResponse(Category), nil
}

func (p *CategoryUseCase) Create(request *models.CreateCategoryRequest) (*models.CategoryResponse, error) {
	err := p.Validate.Struct(request)
	if err != nil {
		return nil, &models.Error{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	Category := &entity.Category{
		Name:        request.Name,
		Description: request.Description,
	}

	createdCategory, err := p.CategoryRepository.Create(Category)
	if err != nil {
		return nil, &models.Error{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	response := converter.CategoryToResponse(createdCategory)
	return response, nil
}

func (p *CategoryUseCase) Update(request *models.UpdateCategoryRequest) (*models.CategoryResponse, error) {

	err := p.Validate.Struct(request)
	if err != nil {
		return nil, &models.Error{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	_, err = p.CategoryRepository.FindByID(request.ID)
	if err != nil {
		return nil, &models.Error{
			Status:  http.StatusNotFound,
			Message: err.Error(),
		}
	}

	Category := &entity.Category{
		ID:          request.ID,
		Name:        request.Name,
		Description: request.Description,
	}

	updatedCategory, err := p.CategoryRepository.Update(Category)
	if err != nil {
		return nil, &models.Error{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	response := converter.CategoryToResponse(updatedCategory)
	return response, nil
}

func (p *CategoryUseCase) Delete(id int) error {
	_, err := p.CategoryRepository.FindByID(id)
	if err != nil {
		return &models.Error{
			Status:  http.StatusNotFound,
			Message: err.Error(),
		}
	}

	err = p.CategoryRepository.Delete(id)
	if err != nil {
		return &models.Error{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil
}
