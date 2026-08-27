package timeline

import (
	"github.com/wyw14/cry-153/internal/model"
	"sync"
	"time"
)

type Window struct {
	mu      sync.RWMutex
	entries map[string][]model.TimelineEntry
}

func NewWindow() *Window { return &Window{entries: map[string][]model.TimelineEntry{}} }
func (w *Window) Add(id, kind, detail string) {
	w.mu.Lock()
	w.entries[id] = append(w.entries[id], model.TimelineEntry{At: time.Now(), Kind: kind, Detail: detail})
	w.mu.Unlock()
}
func (w *Window) Snapshot(id string) []model.TimelineEntry {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return append([]model.TimelineEntry(nil), w.entries[id]...)
}
func (w *Window) SnapshotAll() map[string][]model.TimelineEntry {
	w.mu.RLock()
	defer w.mu.RUnlock()
	out := map[string][]model.TimelineEntry{}
	for id, items := range w.entries {
		out[id] = append([]model.TimelineEntry(nil), items...)
	}
	return out
}
