package main

import (
	"net/http"

	"local.packages/managementapi"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Route("/health", func(r chi.Router) {
		r.Get("/", managementapi.Health)
	})
	r.Route("/products", func(r chi.Router) {
		r.Get("/", managementapi.GetProducts)
	})

	s := &http.Server{
		Addr:    ":3000",
		Handler: r,
	}

	s.ListenAndServe()
}
