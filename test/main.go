package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

type MyError struct {
	Message string `json:"message"`
}

type OuterUrlData struct {
	Url []string `json:"url"`
}

func (err *MyError) Error() string {
	return fmt.Sprintf("error: %s", err.Message)
}

var count int = 0

func handler(w http.ResponseWriter, r *http.Request) {
	file, err := os.ReadFile("./urlData.json")
	if err != nil {
		log.Fatal(err)
	}
	var data OuterUrlData
	err = json.Unmarshal(file, &data)
	if err != nil {
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

	url := data.Url[0]
	endpoint := "/facts/"
	path := chi.URLParam(r, "path")
	animal := r.URL.Query().Get("animal")
	amount := r.URL.Query().Get("amount")

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

	if res.StatusCode != 200 {
		log.Printf("invalid response: %v, status code: %d", string(facts), res.StatusCode)
		http.Error(w, string(facts), res.StatusCode)
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
