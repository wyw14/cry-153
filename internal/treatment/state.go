package treatment

import (
	"github.com/wyw14/cry-153/internal/model"
	"sync"
)

type State struct {
	mu        sync.RWMutex
	applied   map[string]float64
	decisions map[string]model.RecoveryDecision
}

func NewState() *State {
	return &State{applied: map[string]float64{}, decisions: map[string]model.RecoveryDecision{}}
}
func (s *State) Apply(id string, dose float64) { s.mu.Lock(); s.applied[id] = dose; s.mu.Unlock() }
func (s *State) Record(id string, d model.RecoveryDecision) {
	s.mu.Lock()
	s.decisions[id] = d
	s.mu.Unlock()
}
func (s *State) Dose(id string) float64 { s.mu.RLock(); defer s.mu.RUnlock(); return s.applied[id] }
