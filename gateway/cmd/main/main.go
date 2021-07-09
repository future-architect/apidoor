package main

import (
	"net/http"
	"time"

	"local.packages/gateway"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Get("/*", gateway.GetHandler)
		r.Put("/*", gateway.PutHandler)
		r.Delete("/*", gateway.DeleteHandler)
		r.Post("/*", gateway.PostHandler)
	})

	rt := GetEnvWithDeault("READTIMEOUT", 5)
	rht := GetEnvWithDeault("READHEADERTIMEOUT", 5)
	wt := GetEnvWithDeault("WRITETIMEOUT", 20)
	it := GetEnvWithDeault("IDLETIMEOUT", 5)
	mhb := GetEnvWithDeault("MAXHEADERBYTES", 1<<20)

	s := &http.Server{
		Addr:              ":3000",
		Handler:           r,
		ReadTimeout:       time.Duration(rt) * time.Second,
		ReadHeaderTimeout: time.Duration(rht) * time.Second,
		WriteTimeout:      time.Duration(wt) * time.Second,
		IdleTimeout:       time.Duration(it) * time.Second,
		MaxHeaderBytes:    mhb,
	}

	s.ListenAndServe()
}
