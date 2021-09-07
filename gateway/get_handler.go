package gateway

import (
	"errors"
	"log"
	"net/http"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		log.Print("unexpected request content")
		http.Error(w, "unexpected request content", http.StatusBadRequest)
		return
	}

	apikey := r.Header.Get("Authorization")
	if apikey == "" {
		log.Print("No authorization key")
		http.Error(w, "no authorization", http.StatusBadRequest)
		return
	}

	fields, err := DBDriver.GetFields(r.Context(), apikey)
	if err != nil {
		log.Print(err.Error())
		if errors.Is(err, ErrUnauthorizedRequest) {
			http.Error(w, "invalid key or path", http.StatusNotFound)
		} else {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	path, err := fields.URI(r.URL.Path)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "invalid key or path", http.StatusNotFound)
		return
	}

	if err := fields.CheckAPILimit(path); err != nil {
		log.Print(err.Error())
		http.Error(w, "API limit exceeded", http.StatusForbidden)
		return
	}

	query := r.URL.RawQuery
	req, err := http.NewRequest(http.MethodGet, "http://"+path+"?"+query, nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "couldn't make request", http.StatusInternalServerError)
		return
	}
	SetRequestHeader(r, req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("error in http get: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	if err := ResposeChecker(&w, res); err != nil {
		return
	}

	UpdateLog(apikey, path)
}