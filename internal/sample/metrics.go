package sample

import (
	"github.com/wyw14/cry-153/internal/model"
	"sync/atomic"
	"time"
)

type Metrics struct {
	accepted atomic.Uint64
	batches  atomic.Uint64
	last     atomic.Int64
}

func (m *Metrics) Accept(r model.Reading) { m.accepted.Add(1); m.last.Store(r.At.UnixNano()) }
func (m *Metrics) Batch()                 { m.batches.Add(1) }
func (m *Metrics) Accepted() uint64       { return m.accepted.Load() }
func (m *Metrics) Batches() uint64        { return m.batches.Load() }
func (m *Metrics) Last() time.Time {
	v := m.last.Load()
	if v == 0 {
		return time.Time{}
	}
	return time.Unix(0, v)
}
func CopyBatch(batch model.SampleBatch) model.SampleBatch {
	batch.Readings = model.CloneReadings(batch.Readings)
	return batch
}
