package gate

import (
	"sync/atomic"
	"time"
)

type Metrics struct {
	attempts  atomic.Uint64
	successes atomic.Uint64
	failures  atomic.Uint64
	last      atomic.Int64
}

func (m *Metrics) Attempt() { m.attempts.Add(1); m.last.Store(time.Now().UnixNano()) }
func (m *Metrics) Success() { m.successes.Add(1) }
func (m *Metrics) Failure() { m.failures.Add(1) }
func (m *Metrics) Snapshot() map[string]uint64 {
	return map[string]uint64{"attempts": m.attempts.Load(), "successes": m.successes.Load(), "failures": m.failures.Load()}
}
func (m *Metrics) Last() time.Time {
	v := m.last.Load()
	if v == 0 {
		return time.Time{}
	}
	return time.Unix(0, v)
}
