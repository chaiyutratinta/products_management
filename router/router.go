package router

import (
	"net/http"
	"products_management/domain"
	"products_management/middleware"

	"github.com/gorilla/mux"
)

//NewRouter for create router
func NewRouter() *mux.Router {
	router := mux.NewRouter()
	useCase := domain.GetProducts()

	//HandleFunc
	router.HandleFunc("/products", useCase.Get).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/products/{id}", useCase.Delete).Methods(http.MethodDelete, http.MethodGet, http.MethodGet)

	//insert product category
	router.HandleFunc("/category", useCase.AddProductCategory).Methods(http.MethodPost, http.MethodGet)
	router.HandleFunc("/category/{id}", useCase.DeleteProductCategory).Methods(http.MethodDelete)
	router.Use(middleware.AuthMiddleware)

	return router
}
