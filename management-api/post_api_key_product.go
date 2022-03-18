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

// PostAPIKeyProducts godoc
// @Summary Post relationship between api key and authorized products linked to the key
// @Description Post relationship between api key and authorized products linked to the key
// @produce json
// @Param product body model.PostAPIKeyProductsReq true "relationship between apikey and products linked to the apikey"
// @Success 201 {string} string
// @Failure 400 {object} validator.BadRequestResp
// @Failure 500 {string} error
// @Router /keys/products [post]
func PostAPIKeyProducts(w http.ResponseWriter, r *http.Request) {
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

	var req model.PostAPIKeyProductsReq
	if ok := unmarshalJSONAndValidate(w, body.Bytes(), &req); !ok {
		return
	}

	if err := usecase.PostAPIKeyProducts(r.Context(), &req); err != nil {
		writeErrResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "Created")
}
