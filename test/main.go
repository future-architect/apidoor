package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type MyError struct {
	Message string `json:"message"`
}

type OuterUrlData struct {
	Url []string `json:"url"`
}

type OuterKeyData struct {
	Keys map[string][]int `json:"keys"`
}

func (err *MyError) Error() string {
	return fmt.Sprintf("error: %s", err.Message)
}

func contains(list []int, a int) bool {
	for _, v := range list {
		if v == a {
			return true
		}
	}
	return false
}

var count int = 0

func handler(w http.ResponseWriter, r *http.Request) {
	urlfile, err := os.ReadFile("./urlData.json")
	if err != nil {
		log.Fatal(err)
	}
	var urldata OuterUrlData
	if err = json.Unmarshal(urlfile, &urldata); err != nil {
		log.Fatal(err)
	}

	keyfile, err := os.ReadFile("./keyData.json")
	if err != nil {
		log.Fatal(err)
	}
	var keydata OuterKeyData
	if err = json.Unmarshal(keyfile, &keydata); err != nil {
		log.Fatal(err)
	}

	r.Header.Set("Content-Type", "application/json")
	if r.Header.Get("Content-Type") != "application/json" {
		log.Print("unexpected request content")
		http.Error(w, "unexpected request content", http.StatusPreconditionFailed)
		return
	}

	r.Header.Set("Authorization", "testtoken")
	if r.Header.Get("Authorization") == "" {
		log.Print("unauthorized request")
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	apikey := r.URL.Query().Get("apikey")
	apilist, ok := keydata.Keys[apikey]
	if !ok {
		log.Print("api key doesn't exist")
		http.Error(w, "invalid API key", http.StatusForbidden)
		return
	}

	endpoint := "/facts/"
	path := chi.URLParam(r, "path")
	animal := r.URL.Query().Get("animal")
	amount := r.URL.Query().Get("amount")
	apinum, err := strconv.Atoi(r.URL.Query().Get("apinum"))
	if err != nil {
		log.Print("invalid API number")
		http.Error(w, "invalid API number", http.StatusBadRequest)
		return
	}

	if !contains(apilist, apinum) {
		log.Print("unauthorized API request")
		http.Error(w, "unauthorized API request", http.StatusForbidden)
		return
	}
	url := urldata.Url[apinum]

	res, err := http.Get(url + endpoint + path + "?animal_type=" + animal + "&amount=" + amount)
	if err != nil {
		myerr := MyError{
			Message: err.Error(),
		}

		s, err := json.Marshal(myerr)
		if err != nil {
			log.Printf("unexpected error message: %v", err)
			http.Error(w, err.Error(), 500)
			return
		}

		log.Printf("http get: %v", string(s))
		http.Error(w, myerr.Error(), 500)
		return
	}

	facts, err := io.ReadAll(res.Body)
	if err != nil {
		myerr := MyError{
			Message: err.Error(),
		}

		s, err := json.Marshal(myerr)
		if err != nil {
			log.Printf("unexpected error message: %v", err)
			http.Error(w, err.Error(), 500)
			return
		}

		log.Printf("io readall: %v", string(s))
		http.Error(w, myerr.Error(), 500)
		return
	}
	defer res.Body.Close()

	switch code := res.StatusCode; {
	case 400 <= code && code <= 499:
		log.Printf("client error: %v, status code: %d", string(facts), code)
		http.Error(w, string(facts), code)
		return
	case 500 <= code && code <= 599:
		log.Printf("server error: %v, status code: %d", string(facts), code)
		http.Error(w, string(facts), code)
		return
	}

	count++
	log.Printf("called: %d times", count)
	fmt.Fprint(w, string(facts))
}

func main() {
	r := chi.NewRouter()
	r.Route("/{path}", func(r chi.Router) {
		r.Get("/", handler)
		// r.Put("/", putHandler)
		// r.Delete("/", deleteHandler)
		// r.Post("/", postHandler)
	})
	http.ListenAndServe(":3000", r)
}
