package dispatch

import (
	"sync"
	"time"
)

type Lease struct {
	ID    string
	Until time.Time
}
type Leases struct {
	mu    sync.Mutex
	items map[string]Lease
}

func NewLeases() *Leases { return &Leases{items: map[string]Lease{}} }
func (l *Leases) Grant(id string, ttl time.Duration) Lease {
	if ttl <= 0 {
		ttl = time.Minute
	}
	x := Lease{ID: id, Until: time.Now().Add(ttl)}
	l.mu.Lock()
	l.items[id] = x
	l.mu.Unlock()
	return x
}
func (l *Leases) Revoke(id string) { l.mu.Lock(); delete(l.items, id); l.mu.Unlock() }
func (l *Leases) Valid(id string, now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	x, ok := l.items[id]
	return ok && now.Before(x.Until)
}
func (l *Leases) Count() int { l.mu.Lock(); defer l.mu.Unlock(); return len(l.items) }
