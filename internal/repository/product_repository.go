package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"kasir-api/internal/entity"
	"kasir-api/internal/models"
)

type ProductRepository struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{
		DB: db,
	}
}

func (pr *ProductRepository) GetAndCountAll(ctx context.Context, request *models.SearchProductRequest) ([]entity.Product, int, error) {
	query := `
		SELECT id, category_id, name, price, stock
		FROM products
		WHERE 1=1
	`
	countQuery := `
		SELECT COUNT(*)
		FROM products
		WHERE 1=1
	`

	args := make([]interface{}, 0)
	argPos := 1

	if request.Name != "" {
		q := fmt.Sprintf(" AND name ILIKE $%d", argPos)
		query += q
		countQuery += q

		args = append(args, "%"+request.Name+"%")
		argPos++
	}

	if request.CategoryID != "" {
		q := fmt.Sprintf(" AND category_id = $%d", argPos)
		query += q
		countQuery += q

		args = append(args, request.CategoryID)
		argPos++
	}

	exec := executorFromContext(ctx, pr.DB)

	// Count total records
	var total int
	err := exec.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (request.Page - 1) * request.Limit
	query += fmt.Sprintf(" ORDER BY id LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, request.Limit, offset)

	rows, err := exec.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
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
			return nil, 0, err
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return products, total, nil
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

	exec := executorFromContext(ctx, pr.DB)
	row := exec.QueryRowContext(ctx, query, id)

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

	exec := executorFromContext(ctx, pr.DB)
	err := exec.QueryRowContext(
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

	exec := executorFromContext(ctx, pr.DB)
	result, err := exec.ExecContext(
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

	exec := executorFromContext(ctx, pr.DB)
	result, err := exec.ExecContext(ctx, query, id)
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

func (pr *ProductRepository) DeductStockOptimistic(ctx context.Context, id, quantity int, version int) (bool, error) {
	const query = `
		UPDATE products
		SET stock = stock - $1,
		    version = version + 1
		WHERE id = $2 AND version = $3 AND stock >= $1
	`

	exec := executorFromContext(ctx, pr.DB)

	result, err := exec.ExecContext(ctx, query, quantity, id, version)
	if err != nil {
		return false, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return affected == 1, nil
}

func (pr *ProductRepository) FindByIds(ctx context.Context, ids []int) ([]entity.Product, error) {
	if len(ids) == 0 {
		return []entity.Product{}, nil
	}

	query := `
		SELECT *
		FROM products
		WHERE id IN (`

	args := make([]interface{}, 0, len(ids))
	for i, id := range ids {
		if i > 0 {
			query += ","
		}
		query += fmt.Sprintf("$%d", i+1)
		args = append(args, id)
	}
	query += ")"

	exec := executorFromContext(ctx, pr.DB)
	rows, err := exec.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Println("failed to close rows:", err)
		}
	}()

	products := make([]entity.Product, 0)
	for rows.Next() {
		var product entity.Product
		if err := rows.Scan(
			&product.ID,
			&product.CategoryID,
			&product.Name,
			&product.Price,
			&product.Stock,
			&product.Version,
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
