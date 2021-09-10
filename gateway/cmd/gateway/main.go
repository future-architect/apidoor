package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"local.packages/gateway"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Get("/*", gateway.Handler)
		r.Put("/*", gateway.Handler)
		r.Delete("/*", gateway.Handler)
		r.Post("/*", gateway.Handler)
	})

	rt := GetEnvWithDefault("READTIMEOUT", 5)
	rht := GetEnvWithDefault("READHEADERTIMEOUT", 5)
	wt := GetEnvWithDefault("WRITETIMEOUT", 20)
	it := GetEnvWithDefault("IDLETIMEOUT", 5)
	mhb := GetEnvWithDefault("MAXHEADERBYTES", 1<<20)

	s := &http.Server{
		Addr:              ":3000",
		Handler:           r,
		ReadTimeout:       time.Duration(rt) * time.Second,
		ReadHeaderTimeout: time.Duration(rht) * time.Second,
		WriteTimeout:      time.Duration(wt) * time.Second,
		IdleTimeout:       time.Duration(it) * time.Second,
		MaxHeaderBytes:    mhb,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func GetEnvWithDefault(env string, def int) int {
	if def <= 0 {
		log.Fatal("default value should be positive")
	}

	n, err := strconv.Atoi(os.Getenv(env))
	if err != nil {
		log.Fatal(err)
	}
	if n <= 0 {
		n = def
	}

	return n
}
