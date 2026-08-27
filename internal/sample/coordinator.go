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
	closed   bool
}

func NewCoordinator(reg *station.Registry, buffer int) *Coordinator {
	return &Coordinator{registry: reg, output: make(chan model.Reading, buffer)}
}
func (c *Coordinator) Output() <-chan model.Reading { return c.output }
func (c *Coordinator) StartSource(ctx context.Context, source station.Source) {
	c.mu.Lock()
	c.active++
	c.registry.SetOnline(source.ID, true)
	c.mu.Unlock()
	go func() {
		defer func() {
			c.registry.SetOnline(source.ID, false)
			c.mu.Lock()
			c.active--
			if c.active == 0 && !c.closed {
				c.closed = true
				close(c.output)
			}
			c.mu.Unlock()
		}()
		_ = source.Stream(ctx, func(r model.Reading) { c.registry.AddReading(r); c.output <- r })
	}()
}
func (c *Coordinator) Close() {
	c.mu.Lock()
	if !c.closed {
		c.closed = true
		close(c.output)
	}
	c.mu.Unlock()
}
