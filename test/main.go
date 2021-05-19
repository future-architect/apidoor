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
		return
	}

	facts, err := io.ReadAll(res.Body)
	if err != nil {
		myerr := MyError{
			Message: err.Error(),
		}
		log.Printf("io readall: %v", myerr)
		http.Error(w, myerr.Error(), 500)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Printf("invalid response: %v, status code: %d", string(facts), res.StatusCode)
		http.Error(w, string(facts), res.StatusCode)
		return
	}

	fmt.Fprint(w, string(facts))
}

func main() {
	r := chi.NewRouter()
	r.Get("/", handler)
	http.ListenAndServe(":3000", r)
}
