package api

import (
	"encoding/json"
	"net/http"
)

func DecodeJSON(r *http.Request, target any) error { return json.NewDecoder(r.Body).Decode(target) }
func WriteError(w http.ResponseWriter, status int, err error) {
	Encode(w, status, map[string]string{"error": err.Error()})
}
func WriteList(w http.ResponseWriter, key string, items any) {
	Encode(w, http.StatusOK, map[string]any{key: items})
}
