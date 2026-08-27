package treatment

import (
	"sync/atomic"
	"time"
)

type Metrics struct {
	simulations  atomic.Uint64
	applications atomic.Uint64
	last         atomic.Int64
}

func (m *Metrics) Simulated() { m.simulations.Add(1); m.last.Store(time.Now().UnixNano()) }
func (m *Metrics) Applied()   { m.applications.Add(1) }
func (m *Metrics) Snapshot() map[string]uint64 {
	return map[string]uint64{"simulations": m.simulations.Load(), "applications": m.applications.Load()}
}
func (m *Metrics) Last() time.Time {
	v := m.last.Load()
	if v == 0 {
		return time.Time{}
	}
	return time.Unix(0, v)
}
