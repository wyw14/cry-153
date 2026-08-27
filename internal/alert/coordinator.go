package alert

import (
	"context"
	"github.com/wyw14/cry-153/internal/model"
)

type Coordinator struct{ publisher *Publisher }

func NewCoordinator(p *Publisher) *Coordinator { return &Coordinator{publisher: p} }
func (c *Coordinator) Send(ctx context.Context, incident string, intakes []string, active bool) error {
	messages := make(chan Message, len(intakes))
	for _, id := range intakes {
		messages <- Message{IncidentID: incident, IntakeID: id, Active: active, Text: "water protection"}
	}
	close(messages)
	return c.publisher.Stream(ctx, messages)
}
func (c *Coordinator) Alerts() []model.Alert { return c.publisher.Sent() }
