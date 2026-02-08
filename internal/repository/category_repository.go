package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"kasir-api/internal/entity"
)

type CategoryRepository struct {
	DB *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{
		DB: db,
	}
}

func (pr *CategoryRepository) GetAll(ctx context.Context) ([]entity.Category, error) {
	const query = `
		SELECT *
		FROM categories
		ORDER BY id
	`

	exec := executorFromContext(ctx, pr.DB)
	rows, err := exec.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			fmt.Println("failed to close rows:", err)
		}
	}(rows)

	categories := make([]entity.Category, 0)

	for rows.Next() {
		var category entity.Category

		if err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
		); err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (pr *CategoryRepository) FindByID(ctx context.Context, id int) (*entity.Category, error) {
	const query = `
		SELECT *
		FROM categories
		WHERE id = $1
	`

	exec := executorFromContext(ctx, pr.DB)
	row := exec.QueryRowContext(ctx, query, id)

	var category entity.Category

	if err := row.Scan(
		&category.ID,
		&category.Name,
		&category.Description,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	return &category, nil
}

func (pr *CategoryRepository) Create(ctx context.Context, Category *entity.Category) (*entity.Category, error) {
	const query = `
		INSERT INTO categories (name, description)
		VALUES ($1, $2)
		RETURNING id
	`

	exec := executorFromContext(ctx, pr.DB)
	err := exec.QueryRowContext(
		ctx,
		query,
		Category.Name,
		Category.Description,
	).Scan(&Category.ID)
	if err != nil {
		return nil, err
	}

	return Category, nil
}

func (pr *CategoryRepository) Update(ctx context.Context, Category *entity.Category) (*entity.Category, error) {
	const query = `
		UPDATE categories
		SET name = $1, description = $2
		WHERE id = $3
	`

	exec := executorFromContext(ctx, pr.DB)
	result, err := exec.ExecContext(
		ctx,
		query,
		Category.Name,
		Category.Description,
		Category.ID,
	)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, errors.New("category not found")
	}

	return Category, nil
}

func (pr *CategoryRepository) Delete(ctx context.Context, id int) error {
	const query = `
		DELETE FROM categories
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
		return errors.New("category not found")
	}

	return nil
}
