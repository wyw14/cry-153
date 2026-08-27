package incident

import (
	"github.com/wyw14/cry-153/internal/model"
	"sort"
)

func (s *State) List() []model.Incident {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]model.Incident, 0, len(s.items))
	for _, item := range s.items {
		item.Notes = append([]string(nil), item.Notes...)
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.Before(out[j].CreatedAt) })
	return out
}
func Active(items []model.Incident) []model.Incident {
	out := []model.Incident{}
	for _, item := range items {
		if !model.StateTerminal(item.State) {
			out = append(out, item)
		}
	}
	return out
}
func StateCounts(items []model.Incident) map[model.IncidentState]int {
	out := map[model.IncidentState]int{}
	for _, item := range items {
		out[item.State]++
	}
	return out
}
