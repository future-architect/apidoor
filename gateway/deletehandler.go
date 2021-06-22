package gateway

import (
	"io"
	"log"
	"net/http"
)

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		log.Print("unexpected request content")
		http.Error(w, "unexpected request content", http.StatusBadRequest)
		return
	}

	apikey := r.Header.Get("Authorization")
	reqpath := r.URL.Path
	query := r.URL.RawQuery

	path, err := GetAPIURL(r.Context(), apikey, reqpath)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req, err := http.NewRequest(http.MethodDelete, "http://"+path+"?"+query, nil)
	if err != nil {
		log.Print("invalid request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("error in http delete: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	switch code := res.StatusCode; {
	case 400 <= code && code <= 499:
		log.Printf("client error: %v, status code: %d", res.Body, code)
		http.Error(w, "client error", code)
		return
	case 500 <= code && code <= 599:
		log.Printf("server error: %v, status code: %d", res.Body, code)
		http.Error(w, "server error", code)
		return
	}

	UpdateLog(apikey, path)
	if _, err := io.Copy(w, res.Body); err != nil {
		log.Printf("error occur while writing response")
		http.Error(w, "error occur while writing response", http.StatusInternalServerError)
	}
}
