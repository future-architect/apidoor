package model

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type ParamType string

const (
	Header          ParamType = "header"
	Query                     = "query"
	BodyFormEncoded           = "body_form_encoded"
)

type AccessToken struct {
	ParamType ParamType `dynamo:"param_type"`
	Key       string    `dynamo:"key"`
	Value     string    `dynamo:"value"`
}

type NotSupportedParamType string

func (nsp NotSupportedParamType) Error() string { return string(nsp) }

func (at AccessToken) addTokenToRequest(r *http.Request) error {
	switch at.ParamType {
	case Header:
		if r.Header.Get(at.Key) == "" {
			r.Header.Add(at.Key, at.Value)
		}
	case Query:
		if !r.URL.Query().Has(at.Key) {
			query := r.URL.Query()
			query.Add(at.Key, at.Value)
			r.URL.RawQuery = query.Encode()
		}
	case BodyFormEncoded:
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/x-www-form-urlencoded" {
			return fmt.Errorf("content-Type header is not application/x-www-form-urlencoded, got %s", contentType)
		}
		if err := r.ParseForm(); err != nil {
			return fmt.Errorf("reading body as form data failed: %w", err)
		}
		if !r.PostForm.Has(at.Key) {
			r.PostForm.Add(at.Key, at.Value)
			r.Body = io.NopCloser(strings.NewReader(r.PostForm.Encode()))
		}
	default:
		return fmt.Errorf("unsupported param type: %v", at.ParamType)
	}
	return nil
}

type AccessTokens struct {
	Tokens []AccessToken `dynamo:"tokens"`
}

func (ats AccessTokens) AddTokensToRequest(r *http.Request) error {
	var errors error
	for _, v := range ats.Tokens {
		if err := v.addTokenToRequest(r); err != nil {
			if errors == nil {
				errors = err
			} else {
				errors = fmt.Errorf("; %w", err)
			}
		}
	}
	return errors
}
