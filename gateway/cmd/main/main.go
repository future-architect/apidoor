package main

import (
	"net/http"
	"time"

	"local.packages/gateway"

	"github.com/go-chi/chi/v5"
)

func timer() {
	for range time.Tick(time.Minute) {
		gateway.PushLog()
	}
}

func main() {
	go timer()

	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Get("/*", gateway.Handler)
		r.Put("/*", gateway.PutHandler)
		r.Delete("/*", gateway.DeleteHandler)
		r.Post("/*", gateway.PostHandler)
	})
	http.ListenAndServe(":3000", r)
}
