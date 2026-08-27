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

// Acquire consumes a token from the bounded channel.
//
// Invariant: inFlight is incremented only when a token is actually acquired.
// A caller that returns an error has NOT consumed a token and must never call
// Release. Previously the cancelled-wait branch bumped inFlight without holding
// a token, so repeated cancellation of blocked waits inflated the gauge past
// the hard limit (4 -> 7) while the real semaphore stayed sized at 4; the
// desynchronised counter then admitted more concurrent gate close requests
// than the gateway could sustain.
func (c *Capacity) Acquire(ctx context.Context) error {
	select {
	case c.tokens <- struct{}{}:
		c.mu.Lock()
		c.inFlight++
		c.mu.Unlock()
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Release returns a previously acquired token. It is a no-op when no token is
// held, so a caller that received an error from Acquire (and therefore must
// not release) cannot steal a token owned by another in-flight caller.
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
