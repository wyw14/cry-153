package timeline

import (
	"github.com/wyw14/cry-153/internal/model"
	"sort"
	"time"
)

func Windowed(entries []model.TimelineEntry, since time.Time) []model.TimelineEntry {
	out := []model.TimelineEntry{}
	for _, e := range entries {
		if !e.At.Before(since) {
			out = append(out, e)
		}
	}
	return out
}
func Kinds(entries []model.TimelineEntry) []string {
	seen := map[string]bool{}
	for _, e := range entries {
		seen[e.Kind] = true
	}
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
func Latest(entries []model.TimelineEntry) (model.TimelineEntry, bool) {
	if len(entries) == 0 {
		return model.TimelineEntry{}, false
	}
	out := entries[0]
	for _, e := range entries[1:] {
		if e.At.After(out.At) {
			out = e
		}
	}
	return out, true
}
