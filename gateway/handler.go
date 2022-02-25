package gateway

import (
	"context"
	"errors"
	"fmt"
	"github.com/future-architect/apidoor/gateway/datasource"
	"github.com/future-architect/apidoor/gateway/logger"
	"github.com/future-architect/apidoor/gateway/model"
	"io"
	"log"
	"net/http"
	"strings"
)

type DefaultHandler struct {
	Appender   logger.Appender
	DataSource datasource.DataSource
}

func (h DefaultHandler) Handle(w http.ResponseWriter, r *http.Request) {

	apikey := r.Header.Get("X-Apidoor-Authorization")
	if apikey == "" {
		log.Print("No authorization key")
		http.Error(w, "gateway error: no authorization request header", http.StatusBadRequest)
		return
	}

	// get all apis linked with the api key
	fields, err := h.DataSource.GetFields(r.Context(), apikey)
	if err != nil {
		log.Print(err.Error())
		if errors.Is(err, model.ErrUnauthorizedRequest) {
			http.Error(w, "gateway error: invalid key or path", http.StatusNotFound)
		} else {
			http.Error(w, "gateway error: internal server error", http.StatusInternalServerError)
		}
		return
	}

	// look up and check the path
	result, err := fields.LookupTemplate(r.URL.Path)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "gateway error: invalid key or path", http.StatusNotFound)
		return
	}
	forwardURL := result.ForwardURL

	// check if number of request does not exceed limit
	if err := fields.CheckAPILimit(result.Field.Path.JoinPath()); err != nil {
		log.Print(err.Error())
		http.Error(w, "gateway error: API limit exceeded", http.StatusForbidden)
		return
	}

	var req *http.Request
	method := r.Method
	if err := setStoredTokens(r.Context(), result.TemplatePath, r, h.DataSource); err != nil {
		log.Printf("set stored tokens failed: %v", err)
	}

	if method == http.MethodGet || method == http.MethodHead || method == http.MethodDelete || method == http.MethodOptions {
		if r.URL.RawQuery != "" {
			forwardURL = forwardURL + "?" + r.URL.RawQuery
		}
		req, err = http.NewRequest(method, forwardURL, nil)
	} else {
		// Post, Put, Patchなど
		req, err = http.NewRequest(http.MethodPost, forwardURL, r.Body)
	}
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "gateway error: couldn't make request", http.StatusInternalServerError)
		return
	}
	setRequestHeader(r, req)

	// call a target api
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		//TODO: notify detailed errors to the client when certain errors, such as timeout, occurred
		log.Printf("error in http %s: %s", method, err.Error())
		http.Error(w, "gateway error: server error", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	for key, values := range res.Header {
		valueSerialized := strings.Join(values, ",")
		w.Header().Set(key, valueSerialized)
	}
	w.WriteHeader(res.StatusCode)

	// return response and write log
	if err := copyResponse(w, res); err != nil {
		return
	}

	if err := h.Appender.Do(apikey, result.Field.Path.JoinPath(), r, res, calcBillingStatus); err != nil {
		log.Printf("[ERROR] appender write err: %v\n", err)
	}
}

func setRequestHeader(src, dist *http.Request) {
	dist.Header = src.Header
	dist.Header.Del("X-Apidoor-Authorization")
	dist.Header.Del("Connection")
	dist.Header.Del("Cookie")
}

func setStoredTokens(ctx context.Context, templatePath string, src *http.Request, source datasource.DataSource) error {
	apikey := src.Header.Get("X-Apidoor-Authorization")
	accessTokens, err := source.GetAccessTokens(ctx, apikey, templatePath)
	if err != nil {
		return fmt.Errorf("get access tokens failed: %w", err)
	}
	if accessTokens == nil {
		return nil
	}
	if err = accessTokens.AddTokensToRequest(src); err != nil {
		return fmt.Errorf("adding tokens to request failed: %v", err)
	}
	return nil
}

func copyResponse(w http.ResponseWriter, res *http.Response) error {
	if _, err := io.Copy(w, res.Body); err != nil {
		log.Printf("error occur while writing response: %s", err.Error())
		http.Error(w, "gateway error: error occur while writing response", http.StatusInternalServerError)
		return errors.New("error occur while writing response")
	}
	return nil
}

func calcBillingStatus(resp *http.Response) logger.BillingStatus {
	code := resp.StatusCode
	if code >= 500 && code <= 599 {
		return logger.NotBilling
	}
	return logger.Billing
}
