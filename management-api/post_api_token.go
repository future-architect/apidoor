package managementapi

import (
	"bytes"
	"errors"
	"github.com/future-architect/apidoor/managementapi/model"
	"github.com/future-architect/apidoor/managementapi/usecase"
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
		writeErrResponse(w, usecase.NewClientError(errors.New(`unexpected request Content-Type, it must be "application/json"`)))
		return
	}
	body := new(bytes.Buffer)
	if _, err := io.Copy(body, r.Body); err != nil {
		log.Printf("reading request body failed: %v", err)
		writeErrResponse(w, usecase.NewServerError(errors.New(`server error`)))
		return
	}

	var req model.PostAPITokenReq
	if ok := unmarshalJSONAndValidate(w, body.Bytes(), &req); !ok {
		return
	}

	// check whether api routing exists
	if err := usecase.PostAPIToken(r.Context(), req); err != nil {
		writeErrResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "Created")

}
