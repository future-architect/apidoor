package gateway

import (
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
	if apikey == "" {
		log.Print("No authorization key")
		http.Error(w, "no authorization", http.StatusBadRequest)
	}
	reqpath := r.URL.Path
	query := r.URL.RawQuery

	path, err := GetAPIURL(r.Context(), apikey, reqpath)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "invalid key or path", http.StatusNotFound)
		return
	}

	if err := APILimitChecker(apikey, path); err != nil {
		log.Print(err.Error())
		http.Error(w, "API limit exceeded", http.StatusForbidden)
		return
	}

	req, err := http.NewRequest(http.MethodDelete, "http://"+path+"?"+query, nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "couldn't make request", http.StatusInternalServerError)
		return
	}
	RequestHeaderSetter(r, req)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("error occur in http delete: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	if err := ResposeChecker(&w, res); err != nil {
		return
	}

	UpdateLog(apikey, path)
}
