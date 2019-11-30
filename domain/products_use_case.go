package domain

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strings"

	"products_management/controller"
	"products_management/models"
	"products_management/repository"
	"products_management/utils"
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
	UseCase controller.ProductController
}

//GetProducts for get all products
func GetProducts() ProductUseCase {
	client := repository.GetPostgresSession()
	controller := controller.NewController(client)

	return &productUseCase{
		UseCase: controller,
	}
}

func (p *productUseCase) Get(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	products := p.UseCase.GetAllProduct()
	json, _ := json.Marshal(*products)

	writer.WriteHeader(http.StatusOK)
	writer.Write(json)
}

func (p *productUseCase) Add(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	body := &models.Body{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(body)
	utils.Checker(err)

	getError := make(map[string]string)
	if body.Name == "" {
		getError["name"] = "name is require."
	}
	if body.Exp == "" {
		getError["Exp"] = "Exp is require."
	}
	if body.Amount == 0 {
		getError["Amount"] = "Amount at leate 1."
	}

	if len(getError) > 0 {
		json, _ := json.Marshal(getError)
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write(json)

		return
	}

	id := (uuid.New()).String()
	err = p.UseCase.AddProduct(&models.Products{
		Name:     body.Name,
		Exp:      body.Exp,
		Category: body.Category,
		Amount:   body.Amount,
		ID:       id,
	})

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, err.Error())
	}

	json, _ := json.Marshal(struct {
		ID string `json: "id"`
	}{
		ID: id,
	})

	writer.WriteHeader(http.StatusCreated)
	writer.Write(json)
}

func (p *productUseCase) Delete(writer http.ResponseWriter, req *http.Request) {
	id := strings.TrimPrefix(req.URL.Path, "/products/")
	err := p.UseCase.DeleteProduct(&id)

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

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(getBody)
	utils.Checker(err)
	err = p.UseCase.UpdateProduct(&id, getBody)

	if err != nil {
		log.Fatal(err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	writer.WriteHeader(http.StatusNoContent)
}

func (p *productUseCase) GetDetail(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	id := strings.TrimPrefix(req.URL.Path, "/products/")
	result, err := p.UseCase.GetDetailProduct(&id)
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

	getBody := &struct{
		Name	string	`json: "name"`
	}{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(getBody)
	utils.Checker(err)
	categoryName := (*getBody).Name
	err = p.UseCase.AddProductCategory(&categoryName)

	if err != nil {
		log.Fatal(err)
		writer.WriteHeader(http.StatusNotFound)

		return
	}

	writer.WriteHeader(http.StatusOK)
}

func (p *productUseCase) GetProductCategories(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	results, err := p.UseCase.GetProductCategories()
	json, err := json.Marshal(*results)

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
	err := p.UseCase.DeleteProductCategory(&id)

	if err != nil {
		log.Fatal(err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	writer.WriteHeader(http.StatusNoContent)
}