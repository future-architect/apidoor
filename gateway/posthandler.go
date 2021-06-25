package gateway

import (
	"log"
	"net/http"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		log.Print("unexpected request content")
		http.Error(w, "unexpected request content", http.StatusBadRequest)
		return
	}

	apikey := r.Header.Get("Authorization")
	reqpath := r.URL.Path

	path, err := GetAPIURL(r.Context(), apikey, reqpath)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "invalid key or path", http.StatusBadRequest)
		return
	}

	if err := ApiLimitChecker(apikey, path); err != nil {
		log.Print(err.Error())
		http.Error(w, "API limit exceeded", http.StatusBadRequest)
		return
	}

	req, err := http.NewRequest(http.MethodPost, "http://"+path, r.Body)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	RequestHeaderSetter(r, req)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("error in http post: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	if err := ResposeChecker(&w, res); err != nil {
		return
	}

	UpdateLog(apikey, path)
}
