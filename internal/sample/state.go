package sample

import (
	"github.com/wyw14/cry-153/internal/model"
	"sync"
)

type State struct {
	mu      sync.RWMutex
	batches map[string][]model.SampleBatch
}

func NewState() *State { return &State{batches: map[string][]model.SampleBatch{}} }
func (s *State) Add(batch model.SampleBatch) {
	s.mu.Lock()
	s.batches[batch.IncidentID] = append(s.batches[batch.IncidentID], batch)
	s.mu.Unlock()
}
func (s *State) List(id string) []model.SampleBatch {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]model.SampleBatch(nil), s.batches[id]...)
}
