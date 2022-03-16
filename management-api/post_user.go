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

// PostUser godoc
// @Summary Create a user
// @Description Create a user
// @produce json
// @Param user body model.PostUserReq true "user description"
// @Success 201 {string} string
// @Failure 400 {object} validator.BadRequestResp
// @Failure 500 {string} error
// @Router /users [post]
func PostUser(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		log.Printf("unexpected request content: %s", r.Header.Get("Content-Type"))
		writeErrResponse(w, usecase.NewClientError(errors.New(`unexpected request Content-Type, it must be "application/json"`)))
		return
	}
	body := new(bytes.Buffer)
	io.Copy(body, r.Body)

	var req model.PostUserReq
	if ok := unmarshalJSONAndValidate(w, body.Bytes(), &req); !ok {
		return
	}

	if err := usecase.PostUser(r.Context(), req); err != nil {
		writeErrResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "Created")
}
