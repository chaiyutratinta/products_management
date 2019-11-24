package main

import (
	"log"
	"net/http"
	"products_management/router"
	"time"
)

func main() {
	router := router.NewRouter()
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
