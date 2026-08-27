package timeline

import (
	"github.com/wyw14/cry-153/internal/model"
	"sort"
)

func Merge(entries map[string][]model.TimelineEntry) []model.TimelineEntry {
	keys := make([]string, 0, len(entries))
	for k := range entries {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := []model.TimelineEntry{}
	for _, k := range keys {
		out = append(out, entries[k]...)
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].At.Before(out[j].At) })
	return out
}
