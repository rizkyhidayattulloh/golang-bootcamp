package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"kasir-api/internal/entity"
)

var _products = []entity.Product{
	{ID: 1, Name: "Product A", Price: 10000, Stock: 100},
	{ID: 2, Name: "Product B", Price: 20000, Stock: 50},
	{ID: 3, Name: "Product C", Price: 30000, Stock: 75},
}

type ProductRepository struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{
		DB: db,
	}
}

func (pr *ProductRepository) GetAll(ctx context.Context) ([]entity.Product, error) {
	const query = `
		SELECT *
		FROM products
		ORDER BY id
	`

	rows, err := pr.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			fmt.Println("failed to close rows:", err)
		}
	}(rows)

	products := make([]entity.Product, 0)

	for rows.Next() {
		var product entity.Product

		if err := rows.Scan(
			&product.ID,
			&product.CategoryID,
			&product.Name,
			&product.Price,
			&product.Stock,
		); err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (pr *ProductRepository) FindByID(ctx context.Context, id int) (*entity.Product, error) {
	const query = `
		SELECT
			p.id,
			p.name,
			p.price,
			p.stock,
			p.category_id,
			c.id,
			c.name,
			c.description
		FROM products p
		JOIN categories c ON c.id = p.category_id
		WHERE p.id = $1
	`

	row := pr.DB.QueryRowContext(ctx, query, id)

	var product entity.Product
	var category entity.Category

	if err := row.Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.Stock,
		&product.CategoryID,
		&category.ID,
		&category.Name,
		&category.Description,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	product.Category = &category

	return &product, nil
}

func (pr *ProductRepository) Create(ctx context.Context, product *entity.Product) (*entity.Product, error) {
	const query = `
		INSERT INTO products (category_id, name, price, stock)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err := pr.DB.QueryRowContext(
		ctx,
		query,
		product.CategoryID,
		product.Name,
		product.Price,
		product.Stock,
	).Scan(&product.ID)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (pr *ProductRepository) Update(ctx context.Context, product *entity.Product) (*entity.Product, error) {
	const query = `
		UPDATE products
		SET category_id = $1, name = $2, price = $3, stock = $4
		WHERE id = $5
	`

	result, err := pr.DB.ExecContext(
		ctx,
		query,
		product.CategoryID,
		product.Name,
		product.Price,
		product.Stock,
		product.ID,
	)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, errors.New("product not found")
	}

	return product, nil
}

func (pr *ProductRepository) Delete(ctx context.Context, id int) error {
	const query = `
		DELETE FROM products
		WHERE id = $1
	`

	result, err := pr.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("product not found")
	}

	return nil
}
