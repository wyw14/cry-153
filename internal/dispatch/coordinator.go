package dispatch

import (
	"errors"
	"fmt"
	"time"
)

var ErrNoRoute = errors.New("no dispatch route")

type Route interface {
	ID() string
	ETA() time.Duration
}
type route struct {
	id  string
	eta time.Duration
}

func (r *route) ID() string         { return r.id }
func (r *route) ETA() time.Duration { return r.eta }

type Planner struct{ routes map[string]time.Duration }

func NewPlanner() *Planner                          { return &Planner{routes: map[string]time.Duration{}} }
func (p *Planner) Set(id string, eta time.Duration) { p.routes[id] = eta }
func (p *Planner) Plan(id string) (Route, error) {
	eta, ok := p.routes[id]
	if !ok {
		return nil, ErrNoRoute
	}
	return &route{id: id, eta: eta}, nil
}

type Coordinator struct {
	planner *Planner
	state   *State
}

func NewCoordinator(p *Planner, s *State) *Coordinator { return &Coordinator{planner: p, state: s} }
func (c *Coordinator) Create(id string) (string, error) {
	r, err := c.planner.Plan(id)
	if err != nil {
		return "", err
	}
	c.state.Start(id)
	return fmt.Sprintf("dispatch:%s:%s", id, r.ID()), nil
}
func (c *Coordinator) Cancel(id string) { c.state.Finish(id) }
