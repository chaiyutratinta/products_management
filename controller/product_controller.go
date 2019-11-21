package controller

import (
	"log"
	"products_management/models"
	"products_management/repository"

	"go.mongodb.org/mongo-driver/bson"
)

//ProductController ...
type ProductController interface {
	GetAllProduct() *[]models.Products
	AddProduct(*models.Products) error
	DeleteProduct(*string) error
	UpdateProduct(*string, *models.Body) error
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
		log.Fatal(err)

		return &[]models.Products{}
	}

	return products
}

func (r *productController) AddProduct(product *models.Products) error {
	err := r.Repo.Add(product)

	if err != nil {
		log.Fatal(err)

		return err
	}

	return nil
}

func (r *productController) DeleteProduct(id *string) error {
	err := r.Repo.Delete(id)

	if err != nil {
		log.Fatal(err)

		return err
	}

	return nil
}

func (r *productController) UpdateProduct(id *string, body *models.Body) error {
	fields := &bson.D{}

	if body.Name != "" {
		*fields = append(*fields, bson.E{"name", body.Name})
	}
	if body.Exp != "" {
		*fields = append(*fields, bson.E{"exp", body.Exp})
	}
	if len(body.Category) > 0 {
		*fields = append(*fields, bson.E{"category", body.Category})
	}
	if body.Amount != 0 {
		*fields = append(*fields, bson.E{"amount", body.Amount})
	}

	update := bson.D{{"$set", *fields}}
	err := r.Repo.Update(id, &update)

	if err != nil {
		log.Fatal(err)

		return err
	}

	return nil
}
