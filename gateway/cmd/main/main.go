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
		r.Put("/*", gateway.PutHandler)
		r.Delete("/*", gateway.DeleteHandler)
		r.Post("/*", gateway.PostHandler)
	})
	http.ListenAndServe(":3000", r)
}
