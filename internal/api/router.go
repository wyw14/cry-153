package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/wyw14/cry-153/internal/service"
	"net/http"
)

func NewRouter(runtime *service.Runtime) http.Handler {
	r := chi.NewRouter()
	r.Get("/healthz", runtime.Health)
	r.Get("/api/operations", runtime.Operations)
	r.Get("/api/equipment", runtime.Equipment)
	r.Get("/api/interlocks", runtime.Interlocks)
	r.Get("/api/incidents", runtime.Incidents)
	r.Post("/api/incidents", runtime.CreateIncident)
	r.Post("/api/incidents/{id}/close", runtime.CloseIncident)
	return r
}
