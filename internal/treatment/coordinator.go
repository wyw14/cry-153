package treatment

import (
	"github.com/wyw14/cry-153/internal/model"
	"github.com/wyw14/cry-153/internal/sample"
)

type Coordinator struct {
	controller *Controller
	evaluator  sample.Evaluator
}

func NewCoordinator(c *Controller, e sample.Evaluator) *Coordinator {
	return &Coordinator{controller: c, evaluator: e}
}
func (c *Coordinator) Evaluate(id string, window model.SampleWindow) model.RecoveryDecision {
	d := c.controller.Simulate(id, window.Readings)
	return d
}
func (c *Coordinator) Apply(id string, dose float64) { c.controller.Apply(id, dose) }
