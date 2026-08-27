package station

import (
	"github.com/wyw14/cry-153/internal/model"
	"sync/atomic"
	"time"
)

type Metrics struct {
	received atomic.Uint64
	rejected atomic.Uint64
	last     atomic.Int64
}

func (m *Metrics) Record(r model.Reading) { m.received.Add(1); m.last.Store(r.At.UnixNano()) }
func (m *Metrics) Reject()                { m.rejected.Add(1) }
func (m *Metrics) Received() uint64       { return m.received.Load() }
func (m *Metrics) Rejected() uint64       { return m.rejected.Load() }
func (m *Metrics) LastAt() time.Time {
	v := m.last.Load()
	if v == 0 {
		return time.Time{}
	}
	return time.Unix(0, v)
}
func Summarize(reg *Registry) map[string]int {
	snap := reg.Snapshot()
	out := map[string]int{}
	for id, items := range snap {
		out[id] = len(items)
	}
	return out
}
