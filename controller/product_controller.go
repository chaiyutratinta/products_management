package controller

import (
	"fmt"
	"log"
	"products_management/models"
	"products_management/repository"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

//ProductController ...
type ProductController interface {
	GetAllProduct() *[]models.Products
	AddProduct(*models.Products) error
	DeleteProduct(*string) error
	UpdateProduct(*string, *models.Body) error
	GetDetailProduct(*string) (*models.Products, error)

	//insert product category
	AddProductCategory(*string) error
	GetProductCategories() (*[]map[string]string, error)
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

func (r *productController) GetDetailProduct(id *string) (*models.Products, error) {
	filter := bson.D{{"id", *id}}
	result := &models.Body{}
	err := r.Repo.GetDetail(&filter, result)

	if err != nil {
		log.Fatal(err)

		return nil, err
	}

	return &models.Products{
		Name:     result.Name,
		Exp:      result.Exp,
		Category: result.Category,
		Amount:   result.Amount,
	}, nil
}

func (r *productController) AddProductCategory(categoryName *string) error {
	id := (uuid.New()).String()
	sqlCommand := fmt.Sprintf(`INSERT INTO product_category VALUES('%s', '%s')`, id, *categoryName)
	err := r.Repo.AddProtuctCatgegory(&sqlCommand)

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (r *productController) GetProductCategories() (*[]map[string]string, error) {
	sqlCommand := `SELECT * FROM product_category`
	results, err := r.Repo.GetProductCategories(&sqlCommand)

	if err != nil {
		log.Fatal(err)

		return nil, err
	}

	return results, nil
}
