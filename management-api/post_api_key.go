package managementapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/future-architect/apidoor/managementapi/model"
	"github.com/future-architect/apidoor/managementapi/usecase"
	"io"
	"log"
	"net/http"
)

// PostAPIKey godoc
// @Summary post api key
// @Description post api key used for authentication in apidoor gateway
// @produce json
// @Param api_key body model.PostAPIKeyReq true "api key owner"
// @Success 201 {object} model.PostAPIKeyResp
// @Failure 400 {object} validator.BadRequestResp
// @Failure 500 {string} error
// @Router /keys [post]
func PostAPIKey(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		log.Printf("unexpected request content: %s", r.Header.Get("Content-Type"))
		writeErrResponse(w, usecase.NewClientError(errors.New(`unexpected request Content-Type, it must be "application/json"`)))
		return
	}

	body := new(bytes.Buffer)
	if _, err := io.Copy(body, r.Body); err != nil {
		log.Printf("reading request body failed: %v", err)
		writeErrResponse(w, usecase.NewServerError(errors.New(`server error`)))
		return
	}

	var req model.PostAPIKeyReq
	if ok := unmarshalJSONAndValidate(w, body.Bytes(), &req); !ok {
		return
	}

	resp, err := usecase.PostAPIKey(r.Context(), req)
	if err != nil {
		writeErrResponse(w, err)
		return
	}

	respBody, err := json.Marshal(resp)
	if err != nil {
		log.Printf("create json response error: %v", err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(respBody)
}
