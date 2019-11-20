package main

import (
	"log"
	"net/http"
	"products_management/domain"
	"strings"
)

type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type handler struct {
	domain domain.ProductUseCase
}

func (h handler) ServeHTTP(write http.ResponseWriter, req *http.Request) {
	endpoint := "/products"
	path := req.URL.Path

	if !strings.HasPrefix(path, endpoint) {
		write.WriteHeader(http.StatusNotFound)

		return
	}

	suffix := strings.TrimPrefix(path, endpoint)

	switch req.Method {
	case http.MethodGet:
		if suffix != "" {
			http.NotFound(write, req)
			return
		}
		h.domain.Get(write, req)
	case http.MethodPost:
		if suffix != "" {
			http.NotFound(write, req)
			return
		}
		h.domain.Add(write, req)
		return
	case http.MethodDelete:
		h.domain.Delete(write, req)
	}
}

func main() {
	server := &http.Server{
		Addr: ":8080",
		Handler: handler{
			domain: domain.GetProducts(),
		},
	}

	log.Fatal(server.ListenAndServe())
}
