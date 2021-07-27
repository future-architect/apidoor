package managementapi

import (
	"io"
	"net/http"
)

// Health godoc
// @Summary checks if API works
// @Description checks whether this API works correctly or not
// @Produce plain
// @Success 200 {string} string
// @Router /health [get]
func Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "OK")
}
