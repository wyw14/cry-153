package intake

import (
	"context"
	"fmt"
	"github.com/wyw14/cry-153/internal/model"
)

type Coordinator struct {
	locks *Locks
	gate  GateDriver
}
type GateDriver interface {
	Close(context.Context, string, string) error
}

func NewCoordinator(l *Locks, g GateDriver) *Coordinator { return &Coordinator{locks: l, gate: g} }
func (c *Coordinator) CloseGroup(ctx context.Context, incident string, ids []string) error {
	ordered := c.locks.Ordered(ids)
	for _, id := range ordered {
		c.locks.Lock(id)
		defer c.locks.Unlock(id)
	}
	for _, id := range ordered {
		if err := c.gate.Close(ctx, id, incident); err != nil {
			return fmt.Errorf("%s: %w", id, err)
		}
		c.locks.Set(id, false)
	}
	return nil
}
func (c *Coordinator) State() map[string]model.Intake { return c.locks.Snapshot() }
