package main

import (
	"net/http"

	"local.packages/managementapi"

	"github.com/go-chi/chi/v5"
)

// @title Management API
// @version 1.0
// @description This is an API that manages products.
func main() {
	r := chi.NewRouter()
	r.Route("/health", func(r chi.Router) {
		r.Get("/", managementapi.Health)
	})
	r.Route("/products", func(r chi.Router) {
		r.Get("/", managementapi.GetProducts)
	})

	s := &http.Server{
		Addr:    ":3001",
		Handler: r,
	}

	s.ListenAndServe()
}
