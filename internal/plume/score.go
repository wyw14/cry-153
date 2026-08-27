package plume

import (
	"github.com/wyw14/cry-153/internal/model"
	"math"
)

func Score(readings []model.Reading) float64 {
	if len(readings) == 0 {
		return 0
	}
	sum := 0.0
	for _, r := range readings {
		sum += r.Concentration
	}
	return sum / float64(len(readings))
}
func Severity(score float64) string {
	switch {
	case score >= 3:
		return "critical"
	case score >= 2:
		return "high"
	case score >= 1:
		return "elevated"
	default:
		return "normal"
	}
}
func Decay(score float64, hours float64) float64 {
	if score < 0 {
		score = 0
	}
	return score * math.Exp(-hours/24)
}
func ThresholdCrossed(values []float64, threshold float64) bool {
	for _, v := range values {
		if v >= threshold {
			return true
		}
	}
	return false
}
