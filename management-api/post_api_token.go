package managementapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/future-architect/apidoor/managementapi/apirouting"
	"github.com/future-architect/apidoor/managementapi/model"
	"github.com/future-architect/apidoor/managementapi/validator"
	"io"
	"log"
	"net/http"
)

// PostAPIToken godoc
// @Summary post api tokens for call external api
// @Description post api tokens for calling external api
// @produce json
// @Param tokens body model.PostAPITokenReq true "api token description"
// @Success 201 {string} string
// @Failure 400 {object} validator.BadRequestResp
// @Failure 500 {string} error
// @Router /api/token [post]
func PostAPIToken(w http.ResponseWriter, r *http.Request) {
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

	var req model.PostAPITokenReq
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

	// check whether api routing exists
	cnt, err := apirouting.ApiDBDriver.CountRouting(r.Context(), req.APIKey, req.Path)
	if err != nil {
		log.Printf("count api routings db error: %v", err)
		http.Error(w, "server error", http.StatusInternalServerError)
	}
	if cnt == 0 {
		log.Println("api_key or path is wrong")
		resp := validator.NewBadRequestResp("api_key or path is wrong")
		if err := resp.WriteResp(w); err != nil {
			log.Printf("write bad request response failed: %v", err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}

	if err := apirouting.ApiDBDriver.PostAPIToken(r.Context(), req); err != nil {
		log.Printf("insert api token db error: %v", err)
		http.Error(w, "server error", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "Created")

}
