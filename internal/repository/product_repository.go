package repository

import (
	"errors"
	"kasir-api/internal/entity"
)

var _products = []entity.Product{
	{ID: 1, Name: "Product A", Price: 10000, Stock: 100},
	{ID: 2, Name: "Product B", Price: 20000, Stock: 50},
	{ID: 3, Name: "Product C", Price: 30000, Stock: 75},
}

type ProductRepository struct {
}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{}
}

func (pr *ProductRepository) GetAll() ([]entity.Product, error) {
	return _products, nil
}

func (pr *ProductRepository) FindByID(id int) (*entity.Product, error) {
	for _, p := range _products {
		if p.ID == id {
			return &p, nil
		}
	}

	return nil, errors.New("product not found")
}

func (pr *ProductRepository) Create(product *entity.Product) (*entity.Product, error) {
	product.ID = len(_products) + 1
	_products = append(_products, *product)

	return product, nil
}

func (pr *ProductRepository) Update(product *entity.Product) (*entity.Product, error) {
	for i, p := range _products {
		if p.ID == product.ID {
			_products[i] = *product
			return product, nil
		}
	}

	return nil, errors.New("product not found")
}

func (pr *ProductRepository) Delete(id int) error {
	for i, p := range _products {
		if p.ID == id {
			_products = append(_products[:i], _products[i+1:]...)
			return nil
		}
	}

	return errors.New("product not found")
}
