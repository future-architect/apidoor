package gateway

import (
	"errors"
	"io"
	"log"
	"net/http"
)

func ResposeChecker(w *http.ResponseWriter, res *http.Response) error {
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
