package controller

import (
	"products_management/models"
	"products_management/repository"
)

//ProductController ...
type ProductController interface {
	GetAllProduct() *[]models.Products
	AddProduct(*models.Products) error
	DeleteProduct(*string) error
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

func (r *productController) GetAllProduct() *[]models.Products {
	products, err := r.Repo.GetAll()

	if err != nil {
		return &[]models.Products{}
	}

	return products
}

func (r *productController) AddProduct(product *models.Products) error {
	err := r.Repo.Add(product)

	if err != nil {
		return err
	}

	return nil
}

func (r *productController) DeleteProduct(id *string) error {
	err := r.Repo.Delete(id)

	if err != nil {
		return err
	}

	return nil
}
