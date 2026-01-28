package usecase

import (
	"context"
	"database/sql"
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

func (p *CategoryUseCase) List(ctx context.Context) ([]models.CategoryResponse, error) {
	categories, err := p.CategoryRepository.GetAll(ctx)
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

func (p *CategoryUseCase) Get(ctx context.Context, id int) (*models.CategoryResponse, error) {
	Category, err := p.CategoryRepository.FindByID(ctx, id)
	if err != nil {
		return nil, &models.Error{
			Status:  http.StatusNotFound,
			Message: err.Error(),
		}
	}

	return converter.CategoryToResponse(Category), nil
}

func (p *CategoryUseCase) Create(ctx context.Context, request *models.CreateCategoryRequest) (*models.CategoryResponse, error) {
	err := p.Validate.Struct(request)
	if err != nil {
		return nil, &models.Error{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	Category := &entity.Category{
		Name: request.Name,
	}

	if request.Description != nil {
		desc := *request.Description
		Category.Description = sql.NullString{
			String: desc,
			Valid:  desc != "",
		}
	}

	createdCategory, err := p.CategoryRepository.Create(ctx, Category)
	if err != nil {
		return nil, &models.Error{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	response := converter.CategoryToResponse(createdCategory)
	return response, nil
}

func (p *CategoryUseCase) Update(ctx context.Context, request *models.UpdateCategoryRequest) (*models.CategoryResponse, error) {

	err := p.Validate.Struct(request)
	if err != nil {
		return nil, &models.Error{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	_, err = p.CategoryRepository.FindByID(ctx, request.ID)
	if err != nil {
		return nil, &models.Error{
			Status:  http.StatusNotFound,
			Message: err.Error(),
		}
	}

	Category := &entity.Category{
		ID:   request.ID,
		Name: request.Name,
	}

	if request.Description != nil {
		desc := *request.Description
		Category.Description = sql.NullString{
			String: desc,
			Valid:  desc != "",
		}
	}

	updatedCategory, err := p.CategoryRepository.Update(ctx, Category)
	if err != nil {
		return nil, &models.Error{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	response := converter.CategoryToResponse(updatedCategory)
	return response, nil
}

func (p *CategoryUseCase) Delete(ctx context.Context, id int) error {
	_, err := p.CategoryRepository.FindByID(ctx, id)
	if err != nil {
		return &models.Error{
			Status:  http.StatusNotFound,
			Message: err.Error(),
		}
	}

	err = p.CategoryRepository.Delete(ctx, id)
	if err != nil {
		return &models.Error{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil
}
