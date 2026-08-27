package plume

import (
	"github.com/wyw14/cry-153/internal/model"
	"math"
	"sort"
)

func Forecast(readings []model.Reading, horizon int) []float64 {
	if horizon < 1 {
		return nil
	}
	base := 0.0
	if len(readings) > 0 {
		base = readings[len(readings)-1].Concentration
	}
	out := make([]float64, horizon)
	for i := range out {
		out[i] = base * math.Pow(1.05, float64(i+1))
	}
	return out
}
func Quantiles(values []float64) map[string]float64 {
	cp := append([]float64(nil), values...)
	sort.Float64s(cp)
	if len(cp) == 0 {
		return map[string]float64{}
	}
	return map[string]float64{"p50": cp[len(cp)/2], "p90": cp[int(float64(len(cp)-1)*.9)]}
}
func Risk(readings []model.Reading) string {
	if len(readings) == 0 {
		return "unknown"
	}
	v := readings[len(readings)-1].Concentration
	if v > 2 {
		return "high"
	}
	if v > 1 {
		return "elevated"
	}
	return "normal"
}
