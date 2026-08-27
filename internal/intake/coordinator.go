package intake

import (
	"context"
	"fmt"
	"github.com/wyw14/cry-153/internal/model"
	"time"
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
	// Acquire per-intake locks in a globally consistent order (lexicographic)
	// regardless of the caller-supplied ordering. Two concurrent CloseGroup calls
	// that pass the same set of intakes in different orders (e.g. A=[north,south]
	// and B=[south,north]) would otherwise deadlock: A holds north waiting for
	// south while B holds south waiting for north, stalling both gate-close
	// requests and the downstream alert publish permanently.
	ordered := c.locks.Ordered(ids)
	deferred := make([]string, 0, len(ordered))
	defer func() {
		for i := len(deferred) - 1; i >= 0; i-- {
			c.locks.Unlock(deferred[i])
		}
	}()
	for _, id := range ordered {
		c.locks.Lock(id)
		deferred = append(deferred, id)
		time.Sleep(25 * time.Millisecond)
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
