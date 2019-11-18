package main

import (
	"net/http"
	"product_management/api/v1/domain"
)

func main() {
	productHandler := domain.GetProducts()
	
	http.HandleFunc("/", productHandler.Get)
	http.ListenAndServe(":8080", nil)
}
