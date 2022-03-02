package managementapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/future-architect/apidoor/managementapi/model"
	"github.com/future-architect/apidoor/managementapi/validator"
	"io"
	"log"
	"net/http"
)

// PostAPIInfo godoc
// @Summary Get list of API information
// @Description Get list of APIs and its information
// @produce json
// @Param api_info body PostAPIInfoReq true "api information"
// @Success 201 {string} string
// @Failure 400 {object} BadRequestResp
// @Failure 500 {string} error
// @Router /api [post]
func PostAPIInfo(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		log.Printf("unexpected request content: %s", r.Header.Get("Content-Type"))
		resp := validator.NewBadRequestResp(`unexpected request Content-Type, it must be "application/json"`)
		if err := resp.WriteResp(w); err != nil {
			log.Printf("write bad request response failed: %v", err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}
	body := new(bytes.Buffer)
	io.Copy(body, r.Body)

	var req model.PostAPIInfoReq
	if err := json.Unmarshal(body.Bytes(), &req); err != nil {
		if errors.Is(err, model.UnmarshalJsonErr) {
			log.Printf("failed to parse json body: %v", err)
			resp := validator.NewBadRequestResp(model.UnmarshalJsonErr.Error())
			if err := resp.WriteResp(w); err != nil {
				log.Printf("write bad request response failed: %v", err)
				http.Error(w, "server error", http.StatusInternalServerError)
			}
		} else if ve, ok := err.(validator.ValidationErrors); ok {
			log.Printf("input validation failed:\n%v", err)
			if err = ve.ToBadRequestResp().WriteResp(w); err != nil {
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

	if err := db.postAPIInfo(r.Context(), &req); err != nil {
		log.Printf("db insert api info error: %v", err)
		http.Error(w, "server error", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "Created")
}
