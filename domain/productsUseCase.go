package domain

import (
	"fmt"
	"net/http"
	"products_management/controller"
	"products_management/repository"
)

//ProductUseCase ...
type ProductUseCase interface {
	Get(w http.ResponseWriter, r *http.Request)
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

func (p *productUseCase) Get(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, p.UseCase.GetAllProduct())
}
