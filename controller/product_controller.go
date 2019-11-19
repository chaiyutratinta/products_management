package controller

import (
	"products_management/models"
	"products_management/repository"
)

//ProductController ...
type ProductController interface {
	GetAllProduct() []models.Products
}

//ProductsUseCase ...
type productController struct {
	Repo repository.Client
}

//NewController ...
func NewController(r repository.Client) ProductController {

	return &productController{
		Repo: r,
	}
}

func (r *productController) GetAllProduct() []models.Products {
	products := r.Repo.GetAll()

	return []models.Products{products}
}
