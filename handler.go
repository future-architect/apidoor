package apidoor

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

var count int = 0

func Handler(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
	if r.Header.Get("Content-Type") != "application/json" {
		log.Print("unexpected request content")
		http.Error(w, "unexpected request content", http.StatusBadRequest)
		return
	}

	r.Header.Set("Authorization", "testtoken")
	if r.Header.Get("Authorization") == "" {
		log.Print("unauthorized request")
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	url := "https://cat-fact.herokuapp.com"
	apikey := r.URL.Query().Get("apikey")

	apinum, err := GetAPINum(url)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := RequestChecker(apinum, apikey); err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	endpoint := "/facts/"
	path := chi.URLParam(r, "path")
	animal := r.URL.Query().Get("animal")
	amount := r.URL.Query().Get("amount")

	res, err := http.Get(url + endpoint + path + "?animal_type=" + animal + "&amount=" + amount)
	if err != nil {
		log.Printf("error in http get: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	facts, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("error in io.ReadAll: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
