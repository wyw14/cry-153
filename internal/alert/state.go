package alert

import (
	"github.com/wyw14/cry-153/internal/model"
	"sync"
)

type State struct {
	mu     sync.RWMutex
	alerts map[string]model.Alert
}

func NewState() *State             { return &State{alerts: map[string]model.Alert{}} }
func (s *State) Put(a model.Alert) { s.mu.Lock(); s.alerts[a.ID] = a; s.mu.Unlock() }
func (s *State) List() []model.Alert {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]model.Alert, 0, len(s.alerts))
	for _, a := range s.alerts {
		out = append(out, a)
	}
	return out
}
