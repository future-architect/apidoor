package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type MyError struct {
	Message string `json:"message"`
}

func (err *MyError) Error() string {
	return fmt.Sprintf("error: %s", err.Message)
}

func handler(w http.ResponseWriter, r *http.Request) {
	url := "https://cat-fact.herokuapp.com"
	endpoint := "/facts/random"
	animal := "cat"
	amount := 2

	res, err := http.Get(url + endpoint + "?animal_type=" + animal + "&amount=" + fmt.Sprint(amount))
	if err != nil {
		myerr := MyError{
			Message: err.Error(),
		}
		log.Printf("http get: %v", myerr)
		http.Error(w, myerr.Error(), 500)
	}

	facts, err := io.ReadAll(res.Body)
	if err != nil {
		myerr := MyError{
			Message: err.Error(),
		}
		log.Printf("io readall: %v", myerr)
		http.Error(w, myerr.Error(), 500)
	}
	defer res.Body.Close()

	fmt.Fprint(w, string(facts))
}

func main() {
	r := chi.NewRouter()
	r.Get("/", handler)
	http.ListenAndServe(":3000", r)
}
