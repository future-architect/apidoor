package main

import (
	"net/http"

	"local.packages/apidoor"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Route("/{path}", func(r chi.Router) {
		r.Get("/", apidoor.Handler)
		// r.Put("/", putHandler)
		// r.Delete("/", deleteHandler)
		// r.Post("/", postHandler)
	})
	http.ListenAndServe(":3000", r)
}
