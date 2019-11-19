package main

import (
	"net/http"
	"products_management/domain"
)

func main() {
	http.HandleFunc("/products", router)
	http.ListenAndServe(":8080", nil)
}

func router(write http.ResponseWriter, req *http.Request) {
	productsHadler := domain.GetProducts()

	switch req.Method {
	case "GET":
		productsHadler.Get(write, req)
	case "POST":
		productsHadler.Add(write, req)
	}
}
