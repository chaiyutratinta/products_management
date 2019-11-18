package domain

import (
	"net/http"
	"fmt"
	ctl "product_management/api/v1/controller"
	repo "product_management/api/v1/repository"
)
//ProductUseCase ...
type ProductUseCase interface {
	Get(w http.ResponseWriter, r *http.Request)
}

type productUseCase struct {
	UseCase ctl.ProductController
}

//GetProducts for get all products
func GetProducts() ProductUseCase {
	client 		:= repo.GetDbSession()
	controller	:= ctl.NewController(client)

	return &productUseCase {
		UseCase: controller,
	}
}

func (p *productUseCase) Get(w http.ResponseWriter, r *http.Request) {
	fmt.Print(p.UseCase.GetAllProduct())

	fmt.Fprint(w, "First time to get Dildo")
}
