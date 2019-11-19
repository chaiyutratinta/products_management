package domain

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"

	"products_management/controller"
	"products_management/models"
	"products_management/repository"
	"products_management/utils"
)

//ProductUseCase ...
type ProductUseCase interface {
	Get(w http.ResponseWriter, r *http.Request)
	Add(w http.ResponseWriter, r *http.Request)
}

type productUseCase struct {
	UseCase controller.ProductController
}

//GetProducts for get all products
func GetProducts() ProductUseCase {
	client := repository.GetDbSession()
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
	body := struct {
		Name     string   `json: "name"`
		Exp      string   `json: "expire_date"`
		Category []string `json: "category"`
		Amount   int      `json: amount`
	}{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&body)
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

	if len(getError) != 0 {
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
