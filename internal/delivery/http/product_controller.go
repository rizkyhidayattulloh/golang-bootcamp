package http

import (
	"kasir-api/internal/models"
	"kasir-api/internal/usecase"
	"kasir-api/internal/util"
	"net/http"
)

type ProductController struct {
	ProductUseCase *usecase.ProductUseCase
}

func NewProductController(productUseCase *usecase.ProductUseCase) *ProductController {
	return &ProductController{
		ProductUseCase: productUseCase,
	}
}

func (pc *ProductController) GetProducts(w http.ResponseWriter, r *http.Request) error {
	products, err := pc.ProductUseCase.List()
	if err != nil {
		return err
	}

	util.EncodeJSON(w, http.StatusOK, models.WebResponse[[]models.ProductResponse]{
		Data: products,
	})

	return nil
}

func (pc *ProductController) GetProduct(w http.ResponseWriter, r *http.Request, id int) error {
	product, err := pc.ProductUseCase.Get(id)
	if err != nil {
		return err
	}

	util.EncodeJSON(w, http.StatusOK, models.WebResponse[*models.ProductResponse]{
		Data: product,
	})

	return nil
}

func (pc *ProductController) CreateProduct(w http.ResponseWriter, r *http.Request) error {
	var req *models.CreateProductRequest
	err := util.DecodeJSON(r, &req)
	if err != nil {
		return err
	}

	product, err := pc.ProductUseCase.Create(req)
	if err != nil {
		return err
	}

	util.EncodeJSON(w, http.StatusCreated, models.WebResponse[*models.ProductResponse]{
		Data: product,
	})

	return nil
}

func (pc *ProductController) UpdateProduct(w http.ResponseWriter, r *http.Request, id int) error {
	var req *models.UpdateProductRequest
	err := util.DecodeJSON(r, &req)
	if err != nil {
		return err
	}

	req.ID = id

	product, err := pc.ProductUseCase.Update(req)
	if err != nil {
		return err
	}

	util.EncodeJSON(w, http.StatusOK, models.WebResponse[*models.ProductResponse]{
		Data: product,
	})

	return nil
}

func (pc *ProductController) DeleteProduct(w http.ResponseWriter, r *http.Request, id int) error {
	err := pc.ProductUseCase.Delete(id)
	if err != nil {
		return err
	}

	util.EncodeJSON(w, http.StatusOK, models.WebResponse[map[string]string]{
		Data: map[string]string{
			"message": "Product deleted successfully",
		},
	})

	return nil
}
