package http

import (
	"kasir-api/internal/models"
	"kasir-api/internal/usecase"
	"kasir-api/internal/util"
	"net/http"
)

type TransactionController struct {
	TransactionUseCase *usecase.TransactionUseCase
}

func NewTransactionController(transactionUseCase *usecase.TransactionUseCase) *TransactionController {
	return &TransactionController{
		TransactionUseCase: transactionUseCase,
	}
}

func (tc *TransactionController) Checkout(w http.ResponseWriter, r *http.Request) error {
	var req *models.CheckoutRequest
	if err := util.DecodeJSON(r, &req); err != nil {
		return err
	}

	transaction, err := tc.TransactionUseCase.Checkout(r.Context(), req)
	if err != nil {
		return err
	}

	util.EncodeJSON(w, http.StatusCreated, models.WebResponse[any]{
		Data: transaction,
	})

	return nil
}
