package sample

import (
	"context"
	"github.com/wyw14/cry-153/internal/station"
)

type Receiver struct{ coordinator *Coordinator }

func NewReceiver(c *Coordinator) *Receiver { return &Receiver{coordinator: c} }
func (r *Receiver) Receive(ctx context.Context, sources []struct {
	ID     string
	Values []float64
}) {
	for _, item := range sources {
		r.coordinator.StartSource(ctx, station.Source{ID: item.ID, Values: item.Values})
	}
}
