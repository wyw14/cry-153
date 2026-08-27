package treatment

import (
	"github.com/wyw14/cry-153/internal/model"
	"math"
)

type Profile struct {
	BaseDose float64
	MaxDose  float64
	Ramp     float64
}

func DefaultProfile() Profile { return Profile{BaseDose: .5, MaxDose: 4, Ramp: .2} }
func (p Profile) Dose(readiness model.RecoveryDecision) float64 {
	dose := p.BaseDose
	if !readiness.Ready {
		dose += p.Ramp * readiness.Score
	}
	if dose > p.MaxDose {
		dose = p.MaxDose
	}
	return math.Max(0, dose)
}
func (p Profile) Valid() bool { return p.BaseDose >= 0 && p.MaxDose >= p.BaseDose && p.Ramp >= 0 }
func (p Profile) Clamp(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > p.MaxDose {
		return p.MaxDose
	}
	return v
}
