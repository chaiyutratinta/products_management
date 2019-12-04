package router

import (
	"net/http"
	"products_management/domain"

	"github.com/gorilla/mux"
)

//NewRouter for create router
func NewRouter() *mux.Router {
	router := mux.NewRouter()
	useCase := domain.GetProducts()

	//HandleFunc
	router.HandleFunc("/products", useCase.Get).Methods(http.MethodGet)
	router.HandleFunc("/products", useCase.Add).Methods(http.MethodPost)
	router.HandleFunc("/products/{id}", useCase.Delete).Methods(http.MethodDelete)
	router.HandleFunc("/products/{id}", useCase.Edit).Methods(http.MethodPatch)
	router.HandleFunc("/products/{id}", useCase.GetDetail).Methods(http.MethodGet)

	//insert product category
	router.HandleFunc("/category", useCase.AddProductCategory).Methods(http.MethodPost)
	router.HandleFunc("/category", useCase.GetProductCategories).Methods(http.MethodGet)
	router.HandleFunc("/category/{id}", useCase.DeleteProductCategory).Methods(http.MethodDelete)

	return router
}
