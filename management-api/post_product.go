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

// PostProduct godoc
// @Summary Post a product
// @Description Post an API product
// @produce json
// @Param product body model.PostProductReq true "product definition"
// @Success 201 {string} string
// @Failure 400 {object} validator.BadRequestResp
// @Failure 500 {string} error
// @Router /products [post]
func PostProduct(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		log.Printf("unexpected request content: %s", r.Header.Get("Content-Type"))
		writeErrResponse(w, usecase.NewClientError(errors.New(`unexpected request Content-Type, it must be "application/json"`)))
		return
	}
	body := new(bytes.Buffer)
	io.Copy(body, r.Body)

	var req model.PostProductReq
	if ok := unmarshalJSONAndValidate(w, body.Bytes(), &req); !ok {
		return
	}
	req = req.Convert()

	if err := usecase.PostProduct(r.Context(), req); err != nil {
		writeErrResponse(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "Created")
}
