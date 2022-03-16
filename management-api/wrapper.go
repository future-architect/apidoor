package managementapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/future-architect/apidoor/managementapi/model"
	"github.com/future-architect/apidoor/managementapi/usecase"
	"github.com/future-architect/apidoor/managementapi/validator"
	"log"
	"net/http"
)

// unmarshalJSONAndValidate unmarshals body bytes into the target struct and validates the struct
// if the returned error is not nil, it writes 4xx or 5xx status and the response body to ResponseWriter
func unmarshalJSONAndValidate(w http.ResponseWriter, data []byte, req json.Unmarshaler) bool {
	if err := json.Unmarshal(data, &req); err != nil {
		if errors.Is(err, model.UnmarshalJsonErr) {
			log.Printf("failed to parse json body: %v", err)
			resp := validator.NewBadRequestResp(model.UnmarshalJsonErr.Error())
			if respBytes, err := json.Marshal(resp); err == nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write(respBytes)
			} else {
				log.Printf("write bad request response failed: %v", err)
				http.Error(w, "server error", http.StatusInternalServerError)
			}
		} else if ve, ok := err.(validator.ValidationErrors); ok {
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
			// unreachable
			log.Printf("invalid body: %v", err)
			http.Error(w, fmt.Sprintf("invalid body"), http.StatusBadRequest)
		}
		return false
	}
	return true
}

func writeErrResponse(w http.ResponseWriter, err error) {
	if err != nil {
		switch err := err.(type) {
		case usecase.ClientError:
			if respBytes, err := createBadRequestRespBytes(err); err == nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write(respBytes)
			} else {
				log.Printf("write bad request response failed: %v", err)
				http.Error(w, "server error", http.StatusInternalServerError)
			}
		case usecase.ServerError:
			http.Error(w, "server error", http.StatusInternalServerError)
		case validator.ValidationErrors:
			if respBytes, err := json.Marshal(err.ToBadRequestResp()); err == nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write(respBytes)
			} else {
				log.Printf("write bad request response failed: %v", err)
				http.Error(w, "server error", http.StatusInternalServerError)
			}
		default:
			// unexpected error occurred
			log.Printf("unexpected error returned: %v", err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}

	}
}

func createBadRequestRespBytes(err error) ([]byte, error) {
	br := validator.BadRequestResp{
		Message: err.Error(),
	}
	return json.Marshal(br)
}
