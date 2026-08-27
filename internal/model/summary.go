package model

import (
	"sort"
	"time"
)

func SortIncidents(items []Incident) []Incident {
	out := append([]Incident(nil), items...)
	sort.SliceStable(out, func(i, j int) bool { return out[i].UpdatedAt.Before(out[j].UpdatedAt) })
	return out
}
func IncidentAge(item Incident, now time.Time) time.Duration {
	if now.Before(item.CreatedAt) {
		return 0
	}
	return now.Sub(item.CreatedAt)
}
func CloneIncident(item Incident) Incident {
	item.Notes = append([]string(nil), item.Notes...)
	return item
}
func IDs(items []Incident) []string {
	out := make([]string, len(items))
	for i, item := range items {
		out[i] = item.ID
	}
	return out
}
