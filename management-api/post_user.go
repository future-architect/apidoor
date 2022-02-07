package managementapi

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

// PostUser godoc
// @Summary Create a user
// @Description Create a user
// @produce json
// @Param user body PostUserReq true "user description"
// @Success 201 {string} string
// @Failure 400 {object} BadRequestResp
// @Failure 500 {string} error
// @Router /users [post]
func PostUser(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		log.Printf("unexpected request content: %s", r.Header.Get("Content-Type"))
		resp := NewBadRequestResp(`unexpected request Content-Type, it must be "application/json"`)
		if err := resp.writeResp(w); err != nil {
			log.Printf("write bad request response failed: %v", err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}

	var req PostUserReq
	if err := Unmarshal(r.Body, &req); err != nil {
		if errors.Is(err, UnmarshalJsonErr) {
			log.Printf("failed to parse json body: %v", err)
			resp := NewBadRequestResp(UnmarshalJsonErr.Error())
			if err := resp.writeResp(w); err != nil {
				log.Printf("write bad request response failed: %v", err)
				http.Error(w, "server error", http.StatusInternalServerError)
			}
		} else if ve, ok := err.(ValidationErrors); ok {
			log.Printf("input validation failed:\n%v", err)
			if err = ve.toBadRequestResp().writeResp(w); err != nil {
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

	if err := db.postUser(r.Context(), &req); err != nil {
		log.Printf("db insert user error: %v", err)
		http.Error(w, "server error", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "Created")

}
