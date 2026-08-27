package incident

import (
	"github.com/wyw14/cry-153/internal/model"
	"strings"
)

func Sources(items []model.Incident) string {
	parts := make([]string, 0, len(items))
	for _, item := range items {
		parts = append(parts, item.Source)
	}
	return strings.Join(parts, ",")
}
func HasState(items []model.Incident, state model.IncidentState) bool {
	for _, item := range items {
		if item.State == state {
			return true
		}
	}
	return false
}
func Closed(items []model.Incident) int {
	n := 0
	for _, item := range items {
		if item.State == model.StateClosed {
			n++
		}
	}
	return n
}
func Open(items []model.Incident) int { return len(items) - Closed(items) }
