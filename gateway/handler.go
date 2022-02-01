package gateway

import (
	"errors"
	"github.com/future-architect/apidoor/gateway/datasource"
	"github.com/future-architect/apidoor/gateway/logger"
	"github.com/future-architect/apidoor/gateway/model"
	"io"
	"log"
	"net/http"
)

type DefaultHandler struct {
	Appender   logger.Appender
	DataSource datasource.DataSource
}

func (h DefaultHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// check header
	if r.Header.Get("Content-Type") != "application/json" {
		log.Print("unexpected request content")
		http.Error(w, "unexpected request content", http.StatusBadRequest)
		return
	}

	apikey := r.Header.Get("X-Apidoor-Authorization")
	if apikey == "" {
		log.Print("No authorization key")
		http.Error(w, "no authorization request header", http.StatusBadRequest)
		return
	}

	// get all apis linked with the api key
	fields, err := h.DataSource.GetFields(r.Context(), apikey)
	if err != nil {
		log.Print(err.Error())
		if errors.Is(err, model.ErrUnauthorizedRequest) {
			http.Error(w, "invalid key or path", http.StatusNotFound)
		} else {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	// look up and check the path
	path, err := fields.URI(r.URL.Path)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "invalid key or path", http.StatusNotFound)
		return
	}

	// check if number of request does not exceed limit
	if err := fields.CheckAPILimit(path); err != nil {
		log.Print(err.Error())
		http.Error(w, "API limit exceeded", http.StatusForbidden)
		return
	}

	var req *http.Request
	method := r.Method
	query := r.URL.RawQuery

	if method == http.MethodGet || method == http.MethodHead || method == http.MethodDelete || method == http.MethodOptions {
		if query != "" {
			path = path + "?" + query
		}
		req, err = http.NewRequest(method, path, nil)
	} else {
		// Post, Put, Patchなど
		req, err = http.NewRequest(http.MethodPost, path, r.Body)
	}
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "couldn't make request", http.StatusInternalServerError)
		return
	}
	setRequestHeader(r, req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("error in http %s: %s", method, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	// return response and write log
	if err := copyResponse(&w, res); err != nil {
		return
	}

	if err := h.Appender.Do(apikey, path, r); err != nil {
		log.Printf("[ERROR] appender write err: %v\n", err)
	}
}

func setRequestHeader(src, dist *http.Request) {
	dist.Header = src.Header
	dist.Header.Del("X-Apidoor-Authorization")
	dist.Header.Del("Connection")
	dist.Header.Del("Cookie")
}

func copyResponse(w *http.ResponseWriter, res *http.Response) error {
	switch code := res.StatusCode; {
	case 400 <= code && code <= 499:
		log.Printf("client error: %v, status code: %d", res.Body, code)
		http.Error(*w, "client error", code)
		return errors.New("client error")
	case 500 <= code && code <= 599:
		log.Printf("server error: %v, status code: %d", res.Body, code)
		http.Error(*w, "server error", code)
		return errors.New("server error")
	}

	if _, err := io.Copy(*w, res.Body); err != nil {
		log.Printf("error occur while writing response: %s", err.Error())
		http.Error(*w, "error occur while writing response", http.StatusInternalServerError)
		return errors.New("error occur while writing response")
	}

	return nil
}
