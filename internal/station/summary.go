package station

import (
	"fmt"
	"github.com/wyw14/cry-153/internal/model"
)

func Describe(id string, items []model.Reading) string {
	return fmt.Sprintf("station %s has %d readings", id, len(items))
}
func Healthy(items []model.Reading, threshold float64) bool {
	if len(items) == 0 {
		return false
	}
	for _, r := range items {
		if r.Concentration > threshold {
			return false
		}
	}
	return true
}
func Split(items []model.Reading, threshold float64) ([]model.Reading, []model.Reading) {
	low, high := []model.Reading{}, []model.Reading{}
	for _, r := range items {
		if r.Concentration <= threshold {
			low = append(low, r)
		} else {
			high = append(high, r)
		}
	}
	return low, high
}
func Limit(items []model.Reading, n int) []model.Reading {
	if n <= 0 {
		return nil
	}
	if len(items) <= n {
		return model.CloneReadings(items)
	}
	return model.CloneReadings(items[len(items)-n:])
}
