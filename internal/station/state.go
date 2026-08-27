package station

import (
	"github.com/wyw14/cry-153/internal/model"
	"sync"
	"time"
)

type Registry struct {
	mu       sync.RWMutex
	readings map[string][]model.Reading
	online   map[string]bool
}

func NewRegistry() *Registry {
	return &Registry{readings: map[string][]model.Reading{}, online: map[string]bool{}}
}
func (r *Registry) SetOnline(id string, value bool) { r.mu.Lock(); r.online[id] = value; r.mu.Unlock() }
func (r *Registry) AddReading(reading model.Reading) {
	r.mu.Lock()
	r.readings[reading.StationID] = append(r.readings[reading.StationID], reading)
	r.mu.Unlock()
}
func (r *Registry) Snapshot() map[string][]model.Reading { return r.readings }
func (r *Registry) Window(id string, keep int) model.SampleWindow {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := r.readings[id]
	if keep > 0 && len(items) > keep {
		items = items[len(items)-keep:]
	}
	cp := append([]model.Reading(nil), items...)
	return model.SampleWindow{StationID: id, Readings: cp}
}
func (r *Registry) Status() map[string]bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := map[string]bool{}
	for k, v := range r.online {
		out[k] = v
	}
	return out
}
func FreshReading(station string, concentration float64) model.Reading {
	return model.Reading{StationID: station, At: time.Now(), Concentration: concentration, Quality: "ok"}
}
