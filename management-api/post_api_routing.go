package managementapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// PostAPIRouting godoc
// @Summary Post API routing
// @Description Post a new API routing
// @Produce json
// @Param api_routing body PostAPIRoutingReq true "routing parameters"
// @Success 201 {string} string
// @Router /api [post]
func PostAPIRouting(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		log.Print("unexpected request content")
		http.Error(w, "unexpected request content", http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("read body error: %v", err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	var req PostAPIRoutingReq
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Printf("failed to parse json body: %v", err)
		http.Error(w, "failed to parse json body", http.StatusBadRequest)
		return
	}

	if err = req.CheckNoEmptyField(); err != nil {
		log.Printf("invalid body: %v", err)
		http.Error(w, fmt.Sprintf("invalid body: %v", err), http.StatusBadRequest)
		return
	}

	if err = ApiDBDriver.PostAPIRouting(r.Context(), req.ApiKey, req.Path, req.ForwardURL); err != nil {
		log.Printf("post api routing db error: %v", err)
		http.Error(w, "server error", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "Created")
}

type PostAPIRoutingReq struct {
	ApiKey     string `json:"api_key"`
	Path       string `json:"path"`
	ForwardURL string `json:"forward_url"`
}

func (pr PostAPIRoutingReq) CheckNoEmptyField() error {
	if pr.ApiKey == "" {
		return errors.New("api_key field required")
	}
	if pr.Path == "" {
		return errors.New("path field required")
	}
	if pr.ForwardURL == "" {
		return errors.New("forward_url field required")
	}
	return nil
}
