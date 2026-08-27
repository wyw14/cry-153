package alert

import (
	"github.com/wyw14/cry-153/internal/model"
	"sort"
	"strings"
)

func Active(items []model.Alert) []model.Alert {
	out := []model.Alert{}
	for _, a := range items {
		if a.Active {
			out = append(out, a)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.Before(out[j].CreatedAt) })
	return out
}
func ForIncident(items []model.Alert, id string) []model.Alert {
	out := []model.Alert{}
	for _, a := range items {
		if a.IncidentID == id {
			out = append(out, a)
		}
	}
	return out
}
func ContainsMessage(items []model.Alert, token string) bool {
	token = strings.ToLower(token)
	for _, a := range items {
		if strings.Contains(strings.ToLower(a.Message), token) {
			return true
		}
	}
	return false
}
