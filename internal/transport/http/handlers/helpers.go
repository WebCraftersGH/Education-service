package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func decodeJSON(r *http.Request, dst any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(dst)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}

func parseIntQuery(value string, def int) int {
	if strings.TrimSpace(value) == "" {
		return def
	}

	n, err := strconv.Atoi(value)
	if err != nil {
		return def
	}

	return n
}
