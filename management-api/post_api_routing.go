package managementapi

import (
	"errors"
	"fmt"
	"github.com/future-architect/apidoor/managementapi/apirouting"
	"io"
	"log"
	"net/http"
)

// PostAPIRouting godoc
// @Summary Post an API routing
// @Description Post a new API routing
// @Produce json
// @Param api_routing body PostAPIRoutingReq true "routing parameters"
// @Success 201 {string} string
// @Failure 400 {object} BadRequestResp
// @Failure 500 {string} error
// @Router /routing [post]
func PostAPIRouting(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		log.Printf("unexpected request content: %s", r.Header.Get("Content-Type"))
		resp := NewBadRequestResp(`unexpected request Content-Type, it must be "application/json"`)
		if err := resp.writeResp(w); err != nil {
			log.Printf("write bad request response failed: %v", err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}

	var req PostAPIRoutingReq
	if err := Unmarshal(r.Body, &req); err != nil {
		if errors.Is(err, UnmarshalJsonErr) {
			log.Printf("failed to parse json body: %v", err)
			resp := NewBadRequestResp(UnmarshalJsonErr.Error())
			if err := resp.writeResp(w); err != nil {
				log.Printf("write bad request response failed: %v", err)
				http.Error(w, "server error", http.StatusInternalServerError)
			}
		} else if ve, ok := err.(ValidationErrors); ok {
			log.Printf("input validation failed:\n%v", err)
			if err = ve.toBadRequestResp().writeResp(w); err != nil {
				log.Printf("write bad request response failed: %v", err)
				http.Error(w, "server error", http.StatusInternalServerError)
			}
		} else {
			// unreachable
			log.Printf("invalid body: %v", err)
			http.Error(w, fmt.Sprintf("invalid body"), http.StatusBadRequest)
		}
		return
	}

	if err := apirouting.ApiDBDriver.PostRouting(r.Context(), req.ApiKey, req.Path, req.ForwardURL); err != nil {
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
