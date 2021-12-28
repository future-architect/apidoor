package managementapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// PostUser godoc
// @Summary Create a user
// @Description Create a user
// @produce json
// @Param product body PostUserReq true "user description"
// @Success 201 {string} string
// @Failure 400 {object} ValidationFailures
// @Failure 500 {string} error
// @Router /users [post]
func PostUser(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		log.Print("unexpected request content")
		resp := NewValidationFailures(`unexpected request Content-Type, it must be "application/json"`)
		if err := resp.writeResp(w); err != nil {
			log.Printf("write bad request response failed: %v", err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}

	var req PostUserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("failed to parse json body: %v", err)
		resp := NewValidationFailures("failed to parse body as json")
		if err := resp.writeResp(w); err != nil {
			log.Printf("write bad request response failed: %v", err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}

	if err := ValidateStruct(req); err != nil {
		if ve, ok := err.(ValidationErrors); ok {
			log.Printf("input validation failed:\n%v", err)
			if err = ve.toValidationFailures().writeResp(w); err != nil {
				log.Printf("write bad request response failed: %v", err)
				http.Error(w, "server error", http.StatusInternalServerError)
			}
		} else {
			log.Printf("invalid body: %v", err)
			http.Error(w, fmt.Sprintf("invalid body"), http.StatusBadRequest)
		}
		return
	}

	if err := db.postUser(r.Context(), &req); err != nil {
		log.Printf("db insert product error: %v", err)
		http.Error(w, "server error", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "Created")

}
