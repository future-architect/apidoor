package managementapi

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// PostProduct godoc
// @Summary Get list of products
// @Description Get list of APIs and its information
// @produce json
// @Param product body PostProductReq true "api information"
// @Success 201 {string} string
// @Failure 400 {string} error
// @Failure 500 {string} error
// @Router /products [post]
func PostProduct(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		log.Print("unexpected request content")
		http.Error(w, "unexpected request content", http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("read body error: %v", err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	var req PostProductReq
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Printf("failed to parse json body: %v", err)
		http.Error(w, "failed to parse json body", http.StatusBadRequest)
		return
	}

	if err = req.CheckNoEmptyField(); err != nil {
		log.Printf("invalid body: %v", err)
		http.Error(w, fmt.Sprintf("invalid body: %v", err), http.StatusBadRequest)
		return
	}

	_, err = DB.NamedExecContext(r.Context(),
		"INSERT INTO apiinfo(name, source, description, thumbnail, swagger_url) VALUES(:name, :source, :description, :thumbnail, :swagger_url)",
		req)
	if err != nil {
		log.Printf("db insert product error: %v", err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "Created")
}
