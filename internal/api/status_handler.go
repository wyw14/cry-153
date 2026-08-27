package api

import (
	"encoding/json"
	"net/http"
)

func Encode(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func Empty(w http.ResponseWriter) {
	Encode(w, http.StatusOK, map[string]any{"items": []any{}, "status": "ready"})
}
