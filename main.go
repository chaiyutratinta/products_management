package main

import (
	"net/http"
	"products_management/domain"
)

func main() {
	productHandler := domain.GetProducts()

	http.HandleFunc("/", productHandler.Get)
	http.ListenAndServe(":8080", nil)
}
