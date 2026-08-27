package sample

import (
	"context"
	"github.com/wyw14/cry-153/internal/model"
	"github.com/wyw14/cry-153/internal/station"
	"sync"
)

type Coordinator struct {
	registry *station.Registry
	mu       sync.Mutex
	active   int
	output   chan model.Reading
	done     chan struct{} // closed by Close to interrupt in-flight sends
	closed   bool          // Close requested; no new sources may start
	closedOK bool          // output channel has been closed (guarded by c.closed+active)
}

func NewCoordinator(reg *station.Registry, buffer int) *Coordinator {
	return &Coordinator{registry: reg, output: make(chan model.Reading, buffer), done: make(chan struct{})}
}
func (c *Coordinator) Output() <-chan model.Reading { return c.output }
func (c *Coordinator) StartSource(ctx context.Context, source station.Source) {
	c.mu.Lock()
	if c.closed {
		// Shutdown in progress; do not start a source that would outlive the
		// shared output channel.
		c.mu.Unlock()
		return
	}
	c.active++
	c.registry.SetOnline(source.ID, true)
	c.mu.Unlock()
	go func() {
		defer func() {
			c.registry.SetOnline(source.ID, false)
			c.mu.Lock()
			c.active--
			// A single upstream station disconnecting must never close the
			// shared output channel: doing so makes every still-online station
			// panic on its next send ("send on closed channel"). The channel is
			// closed only once Close has been requested AND the last sender has
			// exited, so no send can ever target a closed channel.
			c.maybeCloseOutputLocked()
			c.mu.Unlock()
		}()
		_ = source.Stream(ctx, func(r model.Reading) {
			c.registry.AddReading(r)
			select {
			case c.output <- r:
			case <-c.done:
			case <-ctx.Done():
			}
		})
	}()
}

// maybeCloseOutputLocked closes the shared output channel exactly once, but
// only after Close has been requested and every source goroutine has exited.
// The caller must hold c.mu. The channel reference itself is never cleared so
// that Output() always returns a stable channel (closed once drained, after
// which consumers receive the zero-value/ok=false end-of-stream signal).
func (c *Coordinator) maybeCloseOutputLocked() {
	if c.closed && c.active == 0 && !c.closedOK {
		c.closedOK = true
		close(c.output)
	}
}

func (c *Coordinator) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return
	}
	c.closed = true
	close(c.done) // interrupt any in-flight sends so sources can exit promptly.
	c.maybeCloseOutputLocked()
}
