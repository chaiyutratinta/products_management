package controller

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/go-playground/validator"
	"github.com/google/uuid"

	"products_management/constance"
	"products_management/models"
	"products_management/repository"
)

//ProductController ...
type ProductController interface {
	GetAllProduct() *models.ProductResult
	AddProduct(*models.Products) error
	DeleteProduct(*string) error
	UpdateProduct(*string, *models.Body) (*map[string]string, error)
	GetDetailProduct(*string) (*models.ProductDetail, error)

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

func (r *productController) GetAllProduct() *models.ProductResult {
	products, err := r.GetProducts()

	if err != nil {
		log.Println(err)

		return &models.ProductResult{}
	}

	return products
}

func (r *productController) AddProduct(product *models.Products) error {
	err := r.InsertProduct(product)

	if err != nil {
		log.Println(err)

		return err
	}

	return nil
}

func (r *productController) DeleteProduct(id *string) error {
	err := r.Delete(id)

	if err != nil {
		log.Println(err)

		return err

	}

	return nil
}

func (r *productController) UpdateProduct(id *string, body *models.Body) (*map[string]string, error) {
	validateBody := &struct {
		Name     string `validate:"required"`
		Exp      string `validate:"required,len=0|len=6"`
		Category string `validate:"required"`
		Amount   int    `validate:"required,number,min=1"`
		Price    int    `validate:"required,number,min=1"`
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
	requestErrors := make(map[string]string)

	if category != "" && !r.IsCategoryMatch(&category) {
		requestErrors["category"] = "category not match."
	}

	validate := validator.New()
	errors := validate.Struct(validateBody).(validator.ValidationErrors)

	for _, elm := range errors {
		field := elm.Field()
		tag := elm.ActualTag()

		if tag == "required" {
			delete(mapFiledColumn, field)
		} else {
			requestErrors[mapFiledColumn[field]] = constance.RequestErrors[fmt.Sprintf("%s.%s", field, tag)]
		}
	}

	val := reflect.ValueOf(body).Elem()
	var valPair []string
	numField := val.NumField()

	if len(requestErrors) > 0 || len(errors) == numField {

		return &requestErrors, fmt.Errorf("bad request.")
	}

	for i := 0; i < numField; i++ {
		key := val.Type().Field(i).Name
		value := val.Field(i).Interface()

		if v, ok := mapFiledColumn[key]; ok {
			valPair = append(valPair, fmt.Sprintf("%s='%v'", v, value))
		}
	}

	updateStr := strings.Join(valPair, ", ")

	if err := r.Update(id, &updateStr); err != nil {
		log.Println(err)

		return &requestErrors, err
	}

	return &requestErrors, nil
}

func (r *productController) GetDetailProduct(id *string) (*models.ProductDetail, error) {
	result, err := r.GetDetail(id)

	if err != nil {
		log.Println(err)

		return nil, err
	}

	return result, nil
}

func (r *productController) InsertProductCategory(categoryName *string) error {
	id := (uuid.New()).String()
	err := r.InsertCategory(&id, categoryName)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (r *productController) SelectAllProductCategories() ([]map[string]string, error) {
	results, err := r.GetCategories()

	if err != nil {
		log.Println(err)

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
		log.Println(err)

		return err
	}

	return nil
}

func (r *productController) IsCategoryMatch(id *string) bool {
	result := r.IsCategoryExist(id)

	return result
}
