package service

import (
	"context"
	"sync"
	"time"
)

type Lifecycle struct {
	mu      sync.Mutex
	started bool
	stopped bool
	at      time.Time
}

func (l *Lifecycle) Start() {
	l.mu.Lock()
	l.started = true
	l.stopped = false
	l.at = time.Now()
	l.mu.Unlock()
}
func (l *Lifecycle) Stop()         { l.mu.Lock(); l.stopped = true; l.mu.Unlock() }
func (l *Lifecycle) Running() bool { l.mu.Lock(); defer l.mu.Unlock(); return l.started && !l.stopped }
func (l *Lifecycle) Since(now time.Time) time.Duration {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.at.IsZero() {
		return 0
	}
	return now.Sub(l.at)
}
func RunUntil(ctx context.Context, interval time.Duration, fn func() error) error {
	if interval <= 0 {
		interval = time.Second
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		if err := fn(); err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}
