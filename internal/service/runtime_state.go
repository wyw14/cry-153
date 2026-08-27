package service

import (
	"sync/atomic"
	"time"
)

type RuntimeState struct {
	started atomic.Int64
	stopped atomic.Bool
}

func (s *RuntimeState) Start() { s.started.Store(time.Now().UnixNano()); s.stopped.Store(false) }
func (s *RuntimeState) Stop()  { s.stopped.Store(true) }
func (s *RuntimeState) Started() time.Time {
	v := s.started.Load()
	if v == 0 {
		return time.Time{}
	}
	return time.Unix(0, v)
}
func (s *RuntimeState) Stopped() bool { return s.stopped.Load() }
