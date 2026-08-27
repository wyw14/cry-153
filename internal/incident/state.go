package incident

import (
	"context"
	"github.com/wyw14/cry-153/internal/model"
	"sync"
	"time"
)

type State struct {
	mu       sync.RWMutex
	items    map[string]model.Incident
	contexts map[string]context.CancelFunc
}

func NewState() *State {
	return &State{items: map[string]model.Incident{}, contexts: map[string]context.CancelFunc{}}
}
func (s *State) Put(item model.Incident, cancel context.CancelFunc) {
	s.mu.Lock()
	s.items[item.ID] = item
	s.contexts[item.ID] = cancel
	s.mu.Unlock()
}
func (s *State) Get(id string) (model.Incident, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.items[id]
	return item, ok
}
func (s *State) Update(id string, state model.IncidentState) {
	s.mu.Lock()
	item := s.items[id]
	item.State = state
	item.UpdatedAt = now()
	s.items[id] = item
	s.mu.Unlock()
}
func (s *State) Cancel(id string) {
	s.mu.RLock()
	cancel := s.contexts[id]
	s.mu.RUnlock()
	if cancel != nil {
		cancel()
	}
}
func now() time.Time { return time.Now() }
