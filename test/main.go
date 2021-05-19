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

	url := data.Url[0]
	endpoint := "/facts/random"
	animal := "cat"
	amount := 2

	req, err := http.NewRequest("GET", url+endpoint+"?animal_type="+animal+"&amount="+fmt.Sprint(amount), nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "testtoken")

	if req.Header.Get("Authorization") == "" {
		log.Print("unauthorized request")
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	client := new(http.Client)
	res, err := client.Do(req)
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
	r.Get("/", handler)
	http.ListenAndServe(":3000", r)
}
