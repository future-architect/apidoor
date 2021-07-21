package managementapi

import (
	"io"
	"net/http"
)

// Health checks whether this API works correctly.
func Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "OK")
}
