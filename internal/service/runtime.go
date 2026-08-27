package service

import (
	"encoding/json"
	"github.com/wyw14/cry-153/internal/alert"
	"github.com/wyw14/cry-153/internal/incident"
	"github.com/wyw14/cry-153/internal/journal"
	"github.com/wyw14/cry-153/internal/model"
	"github.com/wyw14/cry-153/internal/station"
	"net/http"
	"time"
)

type Runtime struct {
	IncidentCoord *incident.Coordinator
	State         *incident.State
	Stations      *station.Registry
	Alerts        *alert.Publisher
	Journal       *journal.Store
}

func NewRuntime() *Runtime {
	st := incident.NewState()
	j := journal.NewStore("")
	a := alert.NewPublisher(alert.NewState())
	return &Runtime{IncidentCoord: incident.NewCoordinator(st, j, a), State: st, Stations: station.NewRegistry(), Alerts: a, Journal: j}
}
func (rt *Runtime) Health(w http.ResponseWriter, r *http.Request) {
	write(w, map[string]any{"status": "ok", "time": time.Now().UTC()})
}
func (rt *Runtime) Operations(w http.ResponseWriter, r *http.Request) {
	write(w, map[string]any{"operations": rt.Journal.Entries()})
}
func (rt *Runtime) Equipment(w http.ResponseWriter, r *http.Request) {
	write(w, map[string]any{"stations": rt.Stations.Status()})
}
func (rt *Runtime) Interlocks(w http.ResponseWriter, r *http.Request) {
	write(w, map[string]any{"alerts": rt.Alerts.Sent()})
}
func (rt *Runtime) Incidents(w http.ResponseWriter, r *http.Request) {
	write(w, map[string]any{"incidents": []model.Incident{}})
}
func (rt *Runtime) CreateIncident(w http.ResponseWriter, r *http.Request) {
	item := rt.OpenIncident("api")
	write(w, item)
}
func (rt *Runtime) CloseIncident(w http.ResponseWriter, r *http.Request) {
	write(w, map[string]string{"status": "closed"})
}
func (rt *Runtime) OpenIncident(source string) model.Incident { return rt.IncidentCoord.Open(source) }
func write(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
