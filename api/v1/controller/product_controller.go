package ctl

import (
	"products_management/api/v1/models"
	repo "products_management/api/v1/repository"
)

//ProductController ...
type ProductController interface {
	GetAllProduct() []models.Products
}

//ProductsUseCase ...
type productController struct {
	Repo repo.Client
}

//NewController ...
func NewController(r repo.Client) ProductController {

	return &productController{
		Repo: r,
	}
}

func (r *productController) GetAllProduct() []models.Products {
	products := r.Repo.GetAll()

	return []models.Products{products}
}
