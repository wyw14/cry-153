package api

import (
	"net/http"
	"time"
)

func HealthPayload() map[string]any {
	return map[string]any{"status": "ok", "time": time.Now().UTC().Format(time.RFC3339)}
}
func HealthHandler(w http.ResponseWriter, r *http.Request) { Encode(w, http.StatusOK, HealthPayload()) }
func MethodAllowed(method string) bool                     { return method == http.MethodGet || method == http.MethodPost }
func CacheControl(w http.ResponseWriter)                   { w.Header().Set("Cache-Control", "no-store") }
