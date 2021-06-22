package gateway

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func PutHandler(w http.ResponseWriter, r *http.Request) {
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

	req, err := http.NewRequest(http.MethodPut, "http://"+path+"?"+query, r.Body)
	if err != nil {
		log.Print("invalid request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("error in http put: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	contents, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("error in io.ReadAll: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch code := res.StatusCode; {
	case 400 <= code && code <= 499:
		log.Printf("client error: %v, status code: %d", string(contents), code)
		http.Error(w, string(contents), code)
		return
	case 500 <= code && code <= 599:
		log.Printf("server error: %v, status code: %d", string(contents), code)
		http.Error(w, string(contents), code)
		return
	}

	UpdateLog(apikey, path)
	fmt.Fprint(w, string(contents))
}
