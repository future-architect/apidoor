package managementapi

import (
	"github.com/future-architect/apidoor/managementapi/model"
	"github.com/future-architect/apidoor/managementapi/usecase"
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
		writeErrResponse(w, err)
		return
	}

	// delete item
	if err := usecase.DeleteAPIToken(r.Context(), req); err != nil {
		writeErrResponse(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
