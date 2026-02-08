package usecase

import (
	"context"
	"kasir-api/internal/entity"
	"kasir-api/internal/models"
	"kasir-api/internal/repository"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
)

type TransactionUseCase struct {
	TransactionRepository *repository.TransactionRepository
	ProductRepository     *repository.ProductRepository
	Transactor            *repository.Transactor
	Validate              *validator.Validate
}

func NewTransactionUseCase(transactionRepository *repository.TransactionRepository, productRepository *repository.ProductRepository, transactor *repository.Transactor, validate *validator.Validate) *TransactionUseCase {
	return &TransactionUseCase{
		TransactionRepository: transactionRepository,
		ProductRepository:     productRepository,
		Transactor:            transactor,
		Validate:              validate,
	}
}

func (uc *TransactionUseCase) Checkout(ctx context.Context, req *models.CheckoutRequest) (*models.CheckoutResponse, error) {
	err := uc.Validate.Struct(req)
	if err != nil {
		return nil, &models.Error{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	var response *models.CheckoutResponse
	err = uc.Transactor.WithTx(ctx, func(txCtx context.Context) error {
		transactionProducts := make([]entity.TransactionProduct, 0, len(req.Products))
		totalAmount := 0.0

		productIDs := make([]int, 0, len(req.Products))
		reqProductIDs := make(map[int]int, len(req.Products))
		for _, item := range req.Products {
			if _, ok := reqProductIDs[item.ProductID]; ok {
				reqProductIDs[item.ProductID] += item.Quantity
				continue
			}
			reqProductIDs[item.ProductID] = item.Quantity
			productIDs = append(productIDs, item.ProductID)
		}

		productInfo, err := uc.ProductRepository.FindByIds(txCtx, productIDs)
		if err != nil {
			return &models.Error{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			}
		}

		if len(productInfo) != len(reqProductIDs) {
			return &models.Error{
				Status:  http.StatusBadRequest,
				Message: "there is invalid product ID",
			}
		}

		if err := uc.validateStock(productInfo, req); err != nil {
			return err
		}

		for _, info := range productInfo {
			const maxAttempts = 3
			var price float64
			qty := reqProductIDs[info.ID]

			for attempt := 0; attempt < maxAttempts; attempt++ {
				if info.Stock < qty {
					return &models.Error{
						Status:  http.StatusBadRequest,
						Message: "there is some product that has insufficient stock",
					}
				}

				ok, err := uc.ProductRepository.DeductStockOptimistic(txCtx, info.ID, qty, info.Version)
				if err != nil {
					return &models.Error{
						Status:  http.StatusInternalServerError,
						Message: err.Error(),
					}
				}

				if ok {
					price = info.Price
					info.Stock -= qty
					info.Version++
					break
				}

				if attempt == maxAttempts-1 {
					return &models.Error{
						Status:  http.StatusConflict,
						Message: "stock update conflict, please retry",
					}
				}
			}

			transactionProducts = append(transactionProducts, entity.TransactionProduct{
				ProductID:   info.ID,
				ProductName: info.Name,
				Quantity:    qty,
				UnitPrice:   info.Price,
				Subtotal:    info.Price * float64(qty),
			})

			totalAmount += price * float64(qty)
		}

		transaction := &entity.Transaction{
			TotalAmount: totalAmount,
			CreatedAt:   time.Now(),
		}

		transaction.Products = transactionProducts

		transactionID, err := uc.TransactionRepository.Create(txCtx, transaction)
		if err != nil {
			return &models.Error{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			}
		}

		products := make([]models.CheckoutProductDetail, 0, len(transactionProducts))
		for _, tp := range transactionProducts {
			products = append(products, models.CheckoutProductDetail{
				ProductID: tp.ProductID,
				Quantity:  tp.Quantity,
				Price:     tp.UnitPrice,
			})
		}

		response = &models.CheckoutResponse{
			TransactionID: transactionID,
			TotalAmount:   totalAmount,
			Products:      products,
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (uc *TransactionUseCase) validateStock(products []entity.Product, req *models.CheckoutRequest) error {
	var productStockMap = make(map[int]int)
	for _, p := range products {
		productStockMap[p.ID] = p.Stock
	}

	for _, item := range req.Products {
		if stock, ok := productStockMap[item.ProductID]; !ok || stock < item.Quantity {
			return &models.Error{
				Status:  http.StatusBadRequest,
				Message: "there is some product that has insufficient stock",
			}
		}
	}

	return nil
}
