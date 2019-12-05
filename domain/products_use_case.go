package domain

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"products_management/models"
	"products_management/repository"
	"products_management/services"
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
	services.ProductController
}

//GetProducts for get all products
func GetProducts() ProductUseCase {
	client := repository.GetPostgresSession()
	services := services.NewController(client)

	return &productUseCase{services}
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
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	if resErrs, err := p.AddProduct(body); err != nil {
		switch err.Error() {
		case "validate error.":
			{
				log.Println(err)
				resJson, _ := json.Marshal(resErrs)
				writer.WriteHeader(http.StatusBadRequest)
				writer.Write(resJson)
			}
		default:
			{
				log.Println(err)
				writer.WriteHeader(http.StatusInternalServerError)

				return
			}
		}
	}

	writer.WriteHeader(http.StatusOK)

	return
}

func (p *productUseCase) Delete(writer http.ResponseWriter, req *http.Request) {
	id := strings.TrimPrefix(req.URL.Path, "/products/")
	err := p.DeleteProduct(&id)

	if err != nil {
		log.Println(err)
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
	if err := decoder.Decode(getBody); err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	if reqErrs, err := p.UpdateProduct(&id, getBody); err != nil {
		switch err.Error() {
		case "bad request.":
			{
				resJson, _ := json.Marshal(reqErrs)
				log.Println(err)
				writer.WriteHeader(http.StatusBadRequest)
				writer.Write(resJson)

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
		log.Println(err)
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
		log.Println(err)
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
		log.Println(err)
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
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
