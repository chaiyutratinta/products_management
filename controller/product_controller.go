package controller

import (
	"fmt"
	"log"
	"products_management/models"
	"products_management/repository"
	"reflect"
	"strings"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

//ProductController ...
type ProductController interface {
	GetAllProduct() []models.Products
	AddProduct(*models.Products) error
	DeleteProduct(*string) error
	UpdateProduct(*string, *models.Body) error
	GetDetailProduct(*string) (*models.Products, error)

	//insert product category
	InsertProductCategory(*string) error
	SelectAllProductCategories() ([]map[string]string, error)
	RemoveProductCategory(*string) error
	IsCategoryMatch(*string) bool
}

//ProductsUseCase ...
type productController struct {
	repository.DB
}

//NewController ...
func NewController(db repository.DB) ProductController {

	return &productController{db}
}

func (r *productController) GetAllProduct() []models.Products {
	sqlCommand := fmt.Sprintf(`
		SELECT product.id, product.product_name, product.amount, product.expire, product.price, product_category.category_name  
		FROM product
		LEFT JOIN product_category
		ON product.category_id = product_category.id
	`)
	products, err := r.GetProducts(&sqlCommand)

	if err != nil {
		log.Fatal(err)

		return []models.Products{}
	}

	return products
}

func (r *productController) AddProduct(product *models.Products) error {
	err := r.InsertProduct(product)

	if err != nil {
		log.Fatal(err)

		return err
	}

	return nil
}

func (r *productController) DeleteProduct(id *string) error {
	err := r.Delete(id)

	if err != nil {
		log.Fatal(err)

		return err

	}

	return nil
}

func (r *productController) UpdateProduct(id *string, body *models.Body) error {
	validateBody := &struct {
		Name     string `validate:"required"`
		Exp      string `validate:"required,len=6"`
		Category string `validate:"required"`
		Amount   int    `validate:"min=1"`
		Price    int    `validate:"min=1"`
	}{
		Name:     body.Name,
		Exp:      body.Exp,
		Category: body.Category,
		Amount:   body.Amount,
		Price:    body.Price,
	}
	mapFiledColumn := map[string]string{
		"Name":     "product_name",
		"Exp":      "expire",
		"Category": "category_id",
		"Amount":   "amount",
		"Price":    "price",
	}
	category := validateBody.Category

	if category != "" {
		if ok := r.IsCategoryMatch(&category); !ok {

			return fmt.Errorf("category not match.")
		}
	}

	validate := validator.New()
	errors := validate.Struct(validateBody).(validator.ValidationErrors)

	for _, elm := range errors {
		delete(mapFiledColumn, elm.Field())
	}

	val := reflect.ValueOf(body).Elem()
	var valPair []string

	if len(errors) == val.NumField() {

		return fmt.Errorf("bad request.")
	}

	for i := 0; i < val.NumField(); i++ {
		key := val.Type().Field(i).Name
		value := val.Field(i).Interface()

		if v, ok := mapFiledColumn[key]; ok {
			valPair = append(valPair, fmt.Sprintf("%s='%v'", v, value))
		}
	}

	updateStr := strings.Join(valPair, ", ")

	if err := r.Update(id, &updateStr); err != nil {
		log.Fatal(err)

		return err
	}

	return nil
}

func (r *productController) GetDetailProduct(id *string) (*models.Products, error) {
	result, err := r.GetDetail(id)

	if err != nil {
		log.Fatal(err)

		return nil, err
	}

	return result, nil
}

func (r *productController) InsertProductCategory(categoryName *string) error {
	id := (uuid.New()).String()
	err := r.InsertCategory(&id, categoryName)

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (r *productController) SelectAllProductCategories() ([]map[string]string, error) {
	sqlCommand := `SELECT * FROM product_category`
	results, err := r.GetCategories(&sqlCommand)

	if err != nil {
		log.Fatal(err)

		return nil, err
	}

	categories := []map[string]string{}
	for _, elm := range results {
		categories = append(categories, map[string]string{
			"id":   elm.ID,
			"name": elm.Name,
		})
	}

	return categories, nil
}

func (r *productController) RemoveProductCategory(id *string) error {
	err := r.DeleteCategory(id)

	if err != nil {
		log.Fatal(err)

		return err
	}

	return nil
}

func (r *productController) IsCategoryMatch(id *string) bool {
	sqlCommand := fmt.Sprintf(`SELECT EXISTS (SELECT true FROM product_category WHERE id='%s')`, *id)
	result, err := r.QueryOnce(&sqlCommand)

	if err != nil {
		log.Fatal(err)

		return false
	}

	ismatch, _ := result.(bool)

	return ismatch
}
