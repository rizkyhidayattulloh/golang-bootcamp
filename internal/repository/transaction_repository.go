package repository

import (
	"context"
	"database/sql"
	"fmt"
	"kasir-api/internal/entity"
	"time"
)

type TransactionRepository struct {
	DB *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{
		DB: db,
	}
}

type ReportSummary struct {
	TotalRevenue     float64
	TotalTransaction int
	TopProductName   string
	TopProductQty    int
}

func (r *TransactionRepository) GetReportSummary(ctx context.Context, start, end time.Time) (*ReportSummary, error) {
	exec := executorFromContext(ctx, r.DB)

	const totalsQuery = `
		SELECT
			COALESCE(SUM(total_amount), 0),
			COALESCE(COUNT(*), 0)
		FROM transactions
		WHERE created_at >= $1 AND created_at < $2
	`

	var summary ReportSummary
	if err := exec.QueryRowContext(ctx, totalsQuery, start, end).
		Scan(&summary.TotalRevenue, &summary.TotalTransaction); err != nil {
		return nil, err
	}

	const topProductQuery = `
		SELECT
			tp.product_name AS name,
			SUM(tp.quantity) AS qty
		FROM transaction_products tp
		JOIN transactions t ON t.id = tp.transaction_id
		WHERE t.created_at >= $1 AND t.created_at < $2
		GROUP BY tp.product_id, tp.product_name
		ORDER BY qty DESC
		LIMIT 1
	`

	var name sql.NullString
	var qty sql.NullInt64
	if err := exec.QueryRowContext(ctx, topProductQuery, start, end).
		Scan(&name, &qty); err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if name.Valid {
		summary.TopProductName = name.String
	}
	if qty.Valid {
		summary.TopProductQty = int(qty.Int64)
	}

	return &summary, nil
}

func (r *TransactionRepository) Create(ctx context.Context, transaction *entity.Transaction) (int, error) {
	exec := executorFromContext(ctx, r.DB)

	query := `
		INSERT INTO transactions (total_amount, created_at)
		VALUES ($1, $2)
		RETURNING id
	`

	var transactionID int
	err := exec.QueryRowContext(ctx, query, transaction.TotalAmount, transaction.CreatedAt).Scan(&transactionID)
	if err != nil {
		return 0, err
	}

	transactionProductQuery := `
		INSERT INTO transaction_products (transaction_id, product_id, product_name, quantity, unit_price, subtotal)
		VALUES
	`

	values := []interface{}{}
	for i, tp := range transaction.Products {
		if i > 0 {
			transactionProductQuery += ", "
		}
		lenVal := len(values)
		transactionProductQuery += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", lenVal+1, lenVal+2, lenVal+3, lenVal+4, lenVal+5, lenVal+6)
		values = append(values, transactionID, tp.ProductID, tp.ProductName, tp.Quantity, tp.UnitPrice, tp.Subtotal)
	}

	_, err = exec.ExecContext(ctx, transactionProductQuery, values...)
	if err != nil {
		return 0, err
	}

	return transactionID, nil
}
