package handlers

import (
	"encoding/json"
	"net/http"
)

// CONTENT_TYPE_JSON => json content header
const CONTENT_TYPE_JSON = "application/json"

// MSG_INVALID_HTTP_METHOD => invalid http method message
const MSG_INVALID_HTTP_METHOD = "invalid HTTP method"

// BaseHandler => base handler functionality
type BaseHandler struct {
}

func (base *BaseHandler) writeError(w http.ResponseWriter, err string, code int) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": err})
}
