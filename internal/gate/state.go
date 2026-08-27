package gate

import (
	"context"
	"sync"
)

type Capacity struct {
	mu       sync.Mutex
	limit    int
	inFlight int
	tokens   chan struct{}
}

func NewCapacity(limit int) *Capacity {
	if limit < 1 {
		limit = 1
	}
	return &Capacity{limit: limit, tokens: make(chan struct{}, limit)}
}
func (c *Capacity) Acquire(ctx context.Context) error {
	select {
	case c.tokens <- struct{}{}:
		c.mu.Lock()
		c.inFlight++
		c.mu.Unlock()
		return nil
	case <-ctx.Done():
		c.mu.Lock()
		c.inFlight++
		c.mu.Unlock()
		return ctx.Err()
	}
}
func (c *Capacity) Release() {
	select {
	case <-c.tokens:
		c.mu.Lock()
		if c.inFlight > 0 {
			c.inFlight--
		}
		c.mu.Unlock()
	default:
	}
}
func (c *Capacity) InFlight() int { c.mu.Lock(); defer c.mu.Unlock(); return c.inFlight }
func (c *Capacity) Limit() int    { return c.limit }
