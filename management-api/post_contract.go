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

// PostContract godoc
// @Summary Post a product
// @Description Post an API product
// @produce json
// @Param product body model.PostContractReq true "contract definition"
// @Success 201 {string} string
// @Failure 400 {object} validator.BadRequestResp
// @Failure 500 {string} error
// @Router /contract [post]
func PostContract(w http.ResponseWriter, r *http.Request) {
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

	var req model.PostContractReq
	if ok := unmarshalJSONAndValidate(w, body.Bytes(), &req); !ok {
		return
	}

	if err := usecase.PostContract(r.Context(), req); err != nil {
		writeErrResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "Created")
}
