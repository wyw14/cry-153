package treatment

import (
	"github.com/wyw14/cry-153/internal/model"
	"sort"
)

type Controller struct{ state *State }

func NewController(s *State) *Controller { return &Controller{state: s} }
func (c *Controller) Simulate(id string, readings []model.Reading) model.RecoveryDecision {
	sort.Slice(readings, func(i, j int) bool { return readings[i].Concentration < readings[j].Concentration })
	score := readings[len(readings)/2].Concentration
	return model.RecoveryDecision{Ready: score < 1.0, Score: score, Reason: "simulation"}
}
func (c *Controller) Apply(id string, dose float64) { c.state.Apply(id, dose) }
