package gate

import (
	"context"
	"github.com/wyw14/cry-153/internal/model"
)

type Coordinator struct{ client *Client }

func NewCoordinator(c *Client) *Coordinator { return &Coordinator{client: c} }
func (c *Coordinator) Execute(ctx context.Context, cmd model.GateCommand) error {
	return c.client.Close(ctx, cmd.IntakeID, cmd.IncidentID)
}
func (c *Coordinator) Driver() Driver { return Driver{client: c.client} }

type Driver struct{ client *Client }

func (d Driver) Close(ctx context.Context, intake, incident string) error {
	return d.client.Close(ctx, intake, incident)
}
