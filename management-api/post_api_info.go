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

// PostAPIInfo godoc
// @Summary Get list of API information
// @Description Get list of APIs and its information
// @produce json
// @Param api_info body model.PostAPIInfoReq true "api information"
// @Success 201 {string} string
// @Failure 400 {object} validator.BadRequestResp
// @Failure 500 {string} error
// @Router /api [post]
func PostAPIInfo(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		log.Printf("unexpected request content: %s", r.Header.Get("Content-Type"))
		writeErrResponse(w, usecase.ClientError(errors.New(`unexpected request Content-Type, it must be "application/json"`)))
		return
	}
	body := new(bytes.Buffer)
	io.Copy(body, r.Body)

	var req model.PostAPIInfoReq
	if ok := unmarshalJSONAndValidate(w, body.Bytes(), &req); !ok {
		return
	}

	err := usecase.PostAPIInfo(r.Context(), &req)
	if err != nil {
		writeErrResponse(w, err)
	}
	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "Created")
}
