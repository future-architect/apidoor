package main

import (
	"net/http"

	"local.packages/gateway"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Get("/*", gateway.Handler)
		// r.Put("/", putHandler)
		// r.Delete("/", deleteHandler)
		// r.Post("/", postHandler)
	})
	http.ListenAndServe(":3000", r)
}
