package managementapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// PostAPIRouting godoc
// @Summary Post API routing
// @Description Post a new API routing
// @Produce json
// @Param api_routing body PostAPIRoutingReq true "routing parameters"
// @Success 201 {string} string
// @Failure 400 {string} ValidationFailures
// @Failure 500 {string} error
// @Router /api [post]
func PostAPIRouting(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		log.Print("unexpected request content")
		resp := NewValidationFailures(`unexpected request Content-Type, it must be "application/json"`)
		if err := resp.writeResp(w); err != nil {
			log.Printf("write bad request response failed: %v", err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}

	var req PostAPIRoutingReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("failed to parse json body: %v", err)
		resp := NewValidationFailures("failed to parse body as json")
		if err := resp.writeResp(w); err != nil {
			log.Printf("write bad request response failed: %v", err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}

	if err := ValidateStruct(req); err != nil {
		if ve, ok := err.(ValidationErrors); ok {
			log.Printf("input validation failed:\n%v", err)
			if err = ve.toValidationFailures().writeResp(w); err != nil {
				log.Printf("write bad request response failed: %v", err)
				http.Error(w, "server error", http.StatusInternalServerError)
			}
		} else {
			log.Printf("invalid body: %v", err)
			http.Error(w, fmt.Sprintf("invalid body"), http.StatusBadRequest)
		}
		return
	}

	if err := ApiDBDriver.PostRouting(r.Context(), req.ApiKey, req.Path, req.ForwardURL); err != nil {
		log.Printf("post api routing db error: %v", err)
		http.Error(w, "server error", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "Created")
}

type PostAPIRoutingReq struct {
	ApiKey     string `json:"api_key" validate:"required"`
	Path       string `json:"path" validate:"required"`
	ForwardURL string `json:"forward_url" validate:"required,url"`
}
