package main

import (
	"net/http"
	"time"

	"local.packages/gateway"

	"github.com/go-chi/chi/v5"
)

func timer() {
	for range time.Tick(time.Second) {
		gateway.PushLog()
	}
}

func main() {
	go timer()

	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Get("/*", gateway.Handler)
		// r.Put("/", putHandler)
		// r.Delete("/", deleteHandler)
		// r.Post("/", postHandler)
	})
	http.ListenAndServe(":3000", r)
}
