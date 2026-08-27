package api

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/wyw14/cry-153/internal/service"
	"net/http"
)

func DecodeIncident(r *http.Request) (string, error) {
	var body struct {
		Source string `json:"source"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return "", err
	}
	return body.Source, nil
}
func IncidentID(r *http.Request) string { return chi.URLParam(r, "id") }
func Create(runtime *service.Runtime, w http.ResponseWriter, r *http.Request) {
	source, err := DecodeIncident(r)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	item := runtime.OpenIncident(source)
	_ = json.NewEncoder(w).Encode(item)
}
