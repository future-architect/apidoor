package managementapi

import (
	"encoding/json"
	"github.com/future-architect/apidoor/managementapi/apirouting"
	"github.com/future-architect/apidoor/managementapi/model"
	"github.com/future-architect/apidoor/managementapi/validator"
	"log"
	"net/http"
)

// DeleteAPIToken godoc
// @Summary delete api tokens for call external api
// @Description delete api tokens for calling external api
// @Param api_key body string true "target api_key"
// @Param path body string true "target api_key"
// @Success 204 {object} model.EmptyResp
// @Failure 400 {object} validator.BadRequestResp
// @Failure 500 {string} error
// @Router /api/token [delete]
func DeleteAPIToken(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Printf("parse param error: %v", err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	var req model.DeleteAPITokenReq
	if err := model.SchemaDecoder.Decode(&req, r.Form); err != nil {
		log.Printf("parse query param error: %v", err)
		http.Error(w, "failed to parse query parameters", http.StatusBadRequest)
		return
	}

	if err := validator.ValidateStruct(req); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			log.Printf("input validation failed:\n%v", err)
			if respBytes, err := json.Marshal(ve.ToBadRequestResp()); err == nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write(respBytes)
			} else {
				log.Printf("write bad request response failed: %v", err)
				http.Error(w, "server error", http.StatusInternalServerError)
			}
		} else {
			// unreachable code
			log.Printf("unexpected error returned: %v", err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}

	// delete item
	if err := apirouting.ApiDBDriver.DeleteAPIToken(r.Context(), req); err != nil {
		log.Printf("delete api token db error: %v", err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
