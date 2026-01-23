package repository

import (
	"errors"
	"kasir-api/internal/entity"
	"kasir-api/internal/util"
)

var _categories = []entity.Category{
	{ID: 1, Name: "Category A", Description: util.StringPtr("Description Category A")},
	{ID: 2, Name: "Category B", Description: util.StringPtr("Description Category B")},
	{ID: 3, Name: "Category C", Description: util.StringPtr("Description Category C")},
}

type CategoryRepository struct {
}

func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{}
}

func (pr *CategoryRepository) GetAll() ([]entity.Category, error) {
	return _categories, nil
}

func (pr *CategoryRepository) FindByID(id int) (*entity.Category, error) {
	for _, p := range _categories {
		if p.ID == id {
			return &p, nil
		}
	}

	return nil, errors.New("category not found")
}

func (pr *CategoryRepository) Create(Category *entity.Category) (*entity.Category, error) {
	Category.ID = len(_categories) + 1
	_categories = append(_categories, *Category)

	return Category, nil
}

func (pr *CategoryRepository) Update(Category *entity.Category) (*entity.Category, error) {
	for i, p := range _categories {
		if p.ID == Category.ID {
			_categories[i] = *Category
			return Category, nil
		}
	}

	return nil, errors.New("category not found")
}

func (pr *CategoryRepository) Delete(id int) error {
	for i, p := range _categories {
		if p.ID == id {
			_categories = append(_categories[:i], _categories[i+1:]...)
			return nil
		}
	}

	return errors.New("category not found")
}
