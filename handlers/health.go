package handlers

import (
	"io"
	"log"
	"net/http"
)

// HealthHandler => processes all Health HTTP requests
type HealthHandler struct {
	BaseHandler
}

// Handle => verfies and runs all supported Health HTTP methods
func (h *HealthHandler) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.get(w, r)
		break
	default:
		h.writeError(w, "invalid HTTP method", http.StatusMethodNotAllowed)
		return
	}
}

// Get => returns 200 status with text.  Used for monitoring the web service
func (h *HealthHandler) get(w http.ResponseWriter, r *http.Request) {
	log.Println("Calling Health get()")

	// set headers
	w.Header().Set("Content-Type", CONTENT_TYPE_JSON)
	w.WriteHeader(http.StatusOK)

	// write response
	io.WriteString(w, "green")
}
