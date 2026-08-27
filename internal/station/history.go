package station

import (
	"github.com/wyw14/cry-153/internal/model"
	"sort"
	"time"
)

func Between(items []model.Reading, start, end time.Time) []model.Reading {
	out := []model.Reading{}
	for _, r := range items {
		if !r.At.Before(start) && r.At.Before(end) {
			out = append(out, r)
		}
	}
	return out
}
func Concentrations(items []model.Reading) []float64 {
	out := make([]float64, len(items))
	for i, r := range items {
		out[i] = r.Concentration
	}
	return out
}
func Median(items []model.Reading) float64 {
	values := Concentrations(items)
	if len(values) == 0 {
		return 0
	}
	sort.Float64s(values)
	return values[len(values)/2]
}
func LatestAt(items []model.Reading) time.Time {
	if len(items) == 0 {
		return time.Time{}
	}
	return items[len(items)-1].At
}
