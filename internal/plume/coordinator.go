package plume

import (
	"errors"
	"github.com/wyw14/cry-153/internal/dispatch"
	"github.com/wyw14/cry-153/internal/model"
	"time"
)

type Coordinator struct {
	model     *Model
	state     *State
	scheduler *dispatch.Scheduler
}

func NewCoordinator(m *Model, s *State, q *dispatch.Scheduler) *Coordinator {
	return &Coordinator{model: m, state: s, scheduler: q}
}
func (c *Coordinator) Run(id string, readings []model.Reading) error {
	value, err := c.model.Estimate(readings)
	c.state.Save(id, value, err)
	if err != nil {
		if errors.Is(err, ErrTemporaryUpstream) {
			c.scheduler.Schedule(dispatch.Task{ID: id, Priority: 10, Deadline: timeNow().Add(time.Minute), Kind: "retry"})
		}
		return err
	}
	c.scheduler.Schedule(dispatch.Task{ID: id, Priority: 5, Deadline: timeNow().Add(time.Hour), Kind: "forecast"})
	return nil
}
func timeNow() time.Time { return time.Now() }
