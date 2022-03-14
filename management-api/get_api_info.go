package managementapi

import (
	"encoding/json"
	"github.com/future-architect/apidoor/managementapi/model"
	"github.com/future-architect/apidoor/managementapi/usecase"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

// GetAPIInfo godoc
// @Summary Get list of API info.
// @Description Get list of APIs and its information
// @produce json
// @Success 200 {object} model.APIInfoList
// @Router /api [get]
func GetAPIInfo(w http.ResponseWriter, r *http.Request) {

	list, err := usecase.GetAPIInfo(r.Context())
	if err != nil {
		writeErrResponse(w, err)
		return
	}

	res, err := json.Marshal(model.APIInfoList{List: list})
	if err != nil {
		log.Print("error occurs while reading response")
		writeErrResponse(w, usecase.NewServerError(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
