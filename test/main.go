package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func handler(w http.ResponseWriter, r *http.Request) {
	url := "https://cat-fact.herokuapp.com"
	endpoint := "/facts/random"
	animal := "cat"
	amount := 2

	res, err := http.Get(url + endpoint + "?animal_type=" + animal + "&amount=" + fmt.Sprint(amount))
	if err != nil {
		log.Fatal(err)
	}

	facts, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(w, string(facts))
}

func main() {
	r := chi.NewRouter()
	r.Get("/", handler)
	http.ListenAndServe(":3000", r)
}
