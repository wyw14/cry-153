package dispatch

import "sync"

type State struct {
	mu       sync.RWMutex
	running  map[string]bool
	failures map[string]string
}

func NewState() *State                  { return &State{running: map[string]bool{}, failures: map[string]string{}} }
func (s *State) Start(id string)        { s.mu.Lock(); s.running[id] = true; s.mu.Unlock() }
func (s *State) Finish(id string)       { s.mu.Lock(); delete(s.running, id); s.mu.Unlock() }
func (s *State) Fail(id, reason string) { s.mu.Lock(); s.failures[id] = reason; s.mu.Unlock() }
func (s *State) Running(id string) bool { s.mu.RLock(); defer s.mu.RUnlock(); return s.running[id] }
