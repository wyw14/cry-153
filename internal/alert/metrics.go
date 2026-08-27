package alert

import (
	"sync/atomic"
	"time"
)

type Metrics struct {
	sent    atomic.Uint64
	dropped atomic.Uint64
	last    atomic.Int64
}

func (m *Metrics) Sent()    { m.sent.Add(1); m.last.Store(time.Now().UnixNano()) }
func (m *Metrics) Dropped() { m.dropped.Add(1) }
func (m *Metrics) Snapshot() map[string]uint64 {
	return map[string]uint64{"sent": m.sent.Load(), "dropped": m.dropped.Load()}
}
func (m *Metrics) Last() time.Time {
	v := m.last.Load()
	if v == 0 {
		return time.Time{}
	}
	return time.Unix(0, v)
}
