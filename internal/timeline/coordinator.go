package timeline

import (
	"github.com/wyw14/cry-153/internal/model"
	"github.com/wyw14/cry-153/internal/station"
)

type Coordinator struct {
	window   *Window
	stations *station.Registry
}

func NewCoordinator(w *Window, s *station.Registry) *Coordinator {
	return &Coordinator{window: w, stations: s}
}
func (c *Coordinator) Record(id, kind, detail string)                   { c.window.Add(id, kind, detail) }
func (c *Coordinator) StationSnapshot() map[string][]model.Reading      { return c.stations.Snapshot() }
func (c *Coordinator) IncidentTimeline(id string) []model.TimelineEntry { return c.window.Snapshot(id) }
