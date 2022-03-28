package managementapi

import (
	"encoding/json"
	"github.com/future-architect/apidoor/managementapi/model"
	"github.com/future-architect/apidoor/managementapi/usecase"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

// GetProducts godoc
// @Summary Get list of products.
// @Description Get list of API products
// @produce json
// @Success 200 {object} model.ProductList
// @Router /products [get]
func GetProducts(w http.ResponseWriter, r *http.Request) {

	list, err := usecase.GetProducts(r.Context())
	if err != nil {
		writeErrResponse(w, err)
		return
	}

	res, err := json.Marshal(model.ProductList{List: list})
	if err != nil {
		log.Print("error occurs while reading response")
		writeErrResponse(w, usecase.NewServerError(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
