package station

import (
	"github.com/wyw14/cry-153/internal/model"
	"sort"
)

func OrderedStations(snapshot map[string][]model.Reading) []string {
	ids := make([]string, 0, len(snapshot))
	for id := range snapshot {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}
func Latest(snapshot map[string][]model.Reading) []model.Reading {
	ids := OrderedStations(snapshot)
	out := make([]model.Reading, 0, len(ids))
	for _, id := range ids {
		items := snapshot[id]
		if len(items) > 0 {
			out = append(out, items[len(items)-1])
		}
	}
	return out
}
