package managementapi

import (
	"encoding/json"
	"github.com/future-architect/apidoor/managementapi/model"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

// GetAPIInfo godoc
// @Summary Get list of API info.
// @Description Get list of APIs and its information
// @produce json
// @Success 200 {object} APIInfoList
// @Router /api [get]
func GetAPIInfo(w http.ResponseWriter, r *http.Request) {
	list, err := db.getAPIInfo(r.Context())
	if err != nil {
		log.Printf("execute get apiinfo from db error: %v", err)
		http.Error(w, "error occurs in database", http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(model.APIInfoList{List: list})
	if err != nil {
		log.Print("error occurs while reading response")
		http.Error(w, "error occur in database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
