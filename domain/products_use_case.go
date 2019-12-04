package domain

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"products_management/constance"
	"products_management/controller"
	"products_management/models"
	"products_management/repository"
	"products_management/utils"
	"products_management/validation"

	"github.com/google/uuid"
)

//ProductUseCase ...
type ProductUseCase interface {
	Get(http.ResponseWriter, *http.Request)
	Add(http.ResponseWriter, *http.Request)
	Delete(http.ResponseWriter, *http.Request)
	Edit(http.ResponseWriter, *http.Request)
	GetDetail(http.ResponseWriter, *http.Request)

	//product category
	AddProductCategory(http.ResponseWriter, *http.Request)
	GetProductCategories(http.ResponseWriter, *http.Request)
	DeleteProductCategory(http.ResponseWriter, *http.Request)
}

type productUseCase struct {
	controller.ProductController
}

//GetProducts for get all products
func GetProducts() ProductUseCase {
	client := repository.GetPostgresSession()
	controller := controller.NewController(client)

	return &productUseCase{controller}
}

func (p *productUseCase) Get(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	products := p.GetAllProduct()
	json, _ := json.Marshal(products)

	writer.WriteHeader(http.StatusOK)
	writer.Write(json)
}

func (p *productUseCase) Add(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	body := &models.Body{}

	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()

	if err := decoder.Decode(body); err != nil {
		log.Println(err)
	}

	validateBody := &struct {
		Name     string `validate:"required"`
		Exp      string `validate:"required,len=6"`
		Category string `validate:"-"`
		Amount   int    `validate:"min=1"`
		Price    int    `validate:"min=1"`
	}{
		Name:     body.Name,
		Exp:      body.Exp,
		Category: body.Category,
		Amount:   body.Amount,
		Price:    body.Price,
	}

	validator := validation.New(constance.RequestErrors)
	invalidField := validator.Body(validateBody, models.Body{})

	if ok := p.IsCategoryMatch(&body.Category); !ok {
		invalidField["category"] = constance.RequestErrors["Category"]
	}

	if len(invalidField) > 0 {
		result, _ := json.Marshal(invalidField)
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write(result)

		return
	}

	product := &models.Products{
		ID:       (uuid.New()).String(),
		Name:     body.Name,
		Exp:      body.Exp,
		Category: body.Category,
		Amount:   body.Amount,
		Price:    body.Price,
	}

	if err := p.AddProduct(product); err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	writer.WriteHeader(http.StatusOK)

	return
}

func (p *productUseCase) Delete(writer http.ResponseWriter, req *http.Request) {
	id := strings.TrimPrefix(req.URL.Path, "/products/")
	err := p.DeleteProduct(&id)

	if err != nil {
		log.Fatal(err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	writer.WriteHeader(http.StatusNoContent)
}

func (p *productUseCase) Edit(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	id := strings.TrimPrefix(req.URL.Path, "/products/")
	getBody := &models.Body{}
	defer req.Body.Close()

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(getBody)

	if err != nil {
		log.Println(err)
		fmt.Println(err.Error())
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	err = p.UpdateProduct(&id, getBody)

	if err != nil {
		switch err.Error() {
		case "category not match.", "bad request.":
			{
				log.Println(err)
				writer.WriteHeader(http.StatusBadRequest)

				return
			}

		default:
			{
				log.Println(err)
				writer.WriteHeader(http.StatusInternalServerError)

				return
			}
		}
	}

	writer.WriteHeader(http.StatusNoContent)

	return
}

func (p *productUseCase) GetDetail(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	id := strings.TrimPrefix(req.URL.Path, "/products/")
	result, err := p.GetDetailProduct(&id)
	json, _ := json.Marshal(*result)

	if err != nil {
		log.Fatal(err)
		writer.WriteHeader(http.StatusNotFound)

		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write(json)
}

func (p *productUseCase) AddProductCategory(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	getBody := &struct {
		Name string `json:"name"`
	}{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(getBody)
	utils.Checker(err)
	categoryName := (*getBody).Name
	err = p.InsertProductCategory(&categoryName)

	if err != nil {
		log.Fatal(err)
		writer.WriteHeader(http.StatusNotFound)

		return
	}

	writer.WriteHeader(http.StatusOK)
}

func (p *productUseCase) GetProductCategories(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	results, err := p.SelectAllProductCategories()
	json, err := json.Marshal(results)

	if err != nil {
		log.Fatal(err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write(json)
}

func (p *productUseCase) DeleteProductCategory(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	id := strings.TrimPrefix(req.URL.Path, "/category/")
	err := p.RemoveProductCategory(&id)

	if err != nil {
		log.Fatal(err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
