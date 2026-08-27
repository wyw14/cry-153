package plume

import (
	"github.com/wyw14/cry-153/internal/model"
	"sync"
)

type State struct {
	mu        sync.RWMutex
	estimates map[string]float64
	failures  map[string]error
}

func NewState() *State { return &State{estimates: map[string]float64{}, failures: map[string]error{}} }
func (s *State) Save(id string, value float64, err error) {
	s.mu.Lock()
	if err != nil {
		s.failures[id] = err
	} else {
		s.estimates[id] = value
	}
	s.mu.Unlock()
}
func (s *State) Estimate(id string) (float64, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.estimates[id]
	return v, ok
}
func (s *State) Failure(id string) error { s.mu.RLock(); defer s.mu.RUnlock(); return s.failures[id] }
func ReadingValues(readings []model.Reading) []float64 {
	out := make([]float64, len(readings))
	for i, r := range readings {
		out[i] = r.Concentration
	}
	return out
}
