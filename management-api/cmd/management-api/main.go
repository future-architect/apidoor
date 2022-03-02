package main

import (
	"log"
	"net/http"

	"github.com/future-architect/apidoor/managementapi"
	"github.com/go-chi/chi/v5"
)

// @title Management API
// @version 1.0
// @description This is an API that manages products.
//
// @BasePath /mgmt
func main() {
	r := chi.NewRouter()
	r.Route("/mgmt", func(r chi.Router) {
		r.Route("/health", func(r chi.Router) {
			r.Get("/", managementapi.Health)
		})
		r.Route("/routing", func(r chi.Router) {
			r.Post("/", managementapi.PostAPIRouting)
		})
		r.Route("/api", func(r chi.Router) {
			r.Get("/", managementapi.GetAPIInfo)
			r.Post("/", managementapi.PostAPIInfo)
			r.Get("/search", managementapi.SearchAPIInfo)
			r.Post("/token", managementapi.PostAPIToken)
		})
		r.Route("/users", func(r chi.Router) {
			r.Post("/", managementapi.PostUser)
		})
		r.Route("/products", func(r chi.Router) {
			r.Post("/", managementapi.PostProduct)
		})

	})

	s := &http.Server{
		Addr:    ":3001",
		Handler: r,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
