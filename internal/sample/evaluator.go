package sample

import (
	"github.com/wyw14/cry-153/internal/model"
	"time"
)

type Evaluator struct {
	threshold float64
	span      time.Duration
}

func NewEvaluator(threshold float64, span time.Duration) Evaluator {
	return Evaluator{threshold: threshold, span: span}
}
func (e Evaluator) Evaluate(window model.SampleWindow, now time.Time) model.RecoveryDecision {
	return model.WindowReady(window.Readings, e.threshold, e.span, now)
}
func (e Evaluator) Percentile(values []float64, p float64) float64 {
	if len(values) == 0 {
		return 0
	}
	cp := append([]float64(nil), values...)
	for i := 1; i < len(cp); i++ {
		for j := i; j > 0 && cp[j] < cp[j-1]; j-- {
			cp[j], cp[j-1] = cp[j-1], cp[j]
		}
	}
	idx := int(float64(len(cp)-1) * p)
	return cp[idx]
}
