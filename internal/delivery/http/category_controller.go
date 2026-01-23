package http

import (
	"kasir-api/internal/models"
	"kasir-api/internal/usecase"
	"kasir-api/internal/util"
	"net/http"
)

type CategoryController struct {
	CategoryUseCase *usecase.CategoryUseCase
}

func NewCategoryController(CategoryUseCase *usecase.CategoryUseCase) *CategoryController {
	return &CategoryController{
		CategoryUseCase: CategoryUseCase,
	}
}

func (pc *CategoryController) GetCategories(w http.ResponseWriter, r *http.Request) error {
	categories, err := pc.CategoryUseCase.List()
	if err != nil {
		return err
	}

	util.EncodeJSON(w, http.StatusOK, models.WebResponse[[]models.CategoryResponse]{
		Data: categories,
	})

	return nil
}

func (pc *CategoryController) GetCategory(w http.ResponseWriter, r *http.Request, id int) error {
	Category, err := pc.CategoryUseCase.Get(id)
	if err != nil {
		return err
	}

	util.EncodeJSON(w, http.StatusOK, models.WebResponse[*models.CategoryResponse]{
		Data: Category,
	})

	return nil
}

func (pc *CategoryController) CreateCategory(w http.ResponseWriter, r *http.Request) error {
	var req *models.CreateCategoryRequest
	err := util.DecodeJSON(r, &req)
	if err != nil {
		return err
	}

	Category, err := pc.CategoryUseCase.Create(req)
	if err != nil {
		return err
	}

	util.EncodeJSON(w, http.StatusCreated, models.WebResponse[*models.CategoryResponse]{
		Data: Category,
	})

	return nil
}

func (pc *CategoryController) UpdateCategory(w http.ResponseWriter, r *http.Request, id int) error {
	var req *models.UpdateCategoryRequest
	err := util.DecodeJSON(r, &req)
	if err != nil {
		return err
	}

	req.ID = id

	Category, err := pc.CategoryUseCase.Update(req)
	if err != nil {
		return err
	}

	util.EncodeJSON(w, http.StatusOK, models.WebResponse[*models.CategoryResponse]{
		Data: Category,
	})

	return nil
}

func (pc *CategoryController) DeleteCategory(w http.ResponseWriter, r *http.Request, id int) error {
	err := pc.CategoryUseCase.Delete(id)
	if err != nil {
		return err
	}

	util.EncodeJSON(w, http.StatusOK, models.WebResponse[map[string]string]{
		Data: map[string]string{
			"message": "Category deleted successfully",
		},
	})

	return nil
}
