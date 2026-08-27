package incident

import (
	"context"
	"github.com/wyw14/cry-153/internal/alert"
	"github.com/wyw14/cry-153/internal/journal"
	"github.com/wyw14/cry-153/internal/model"
	"time"
)

type Coordinator struct {
	state   *State
	journal *journal.Store
	alerts  *alert.Publisher
}

func NewCoordinator(st *State, j *journal.Store, a *alert.Publisher) *Coordinator {
	return &Coordinator{state: st, journal: j, alerts: a}
}
func (c *Coordinator) Open(source string) model.Incident {
	ctx, cancel := context.WithCancel(context.Background())
	item := model.Incident{ID: model.NewID("incident"), State: model.StateSuspected, CreatedAt: time.Now(), UpdatedAt: time.Now(), Source: source}
	c.state.Put(item, cancel)
	_ = c.journal.Append(journal.Timeline("incident.open", item.ID))
	_ = ctx
	return item
}
func (c *Coordinator) Confirm(id string) error {
	c.state.Update(id, model.StateConfirmed)
	return c.journal.Append(journal.Timeline("incident.confirm", id))
}
func (c *Coordinator) Close(id string) error {
	c.state.Cancel(id)
	c.state.Update(id, model.StateClosed)
	if c.alerts != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_ = c.alerts.PublishContext(ctx, alert.Message{IncidentID: id, Active: false, Text: "incident closed"})
	}
	return c.journal.Append(journal.Timeline("incident.close", id))
}
