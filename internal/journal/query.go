package journal

import (
	"github.com/wyw14/cry-153/internal/model"
	"sort"
	"strings"
	"time"
)

func Filter(entries []model.TimelineEntry, kind string) []model.TimelineEntry {
	out := []model.TimelineEntry{}
	for _, e := range entries {
		if kind == "" || e.Kind == kind {
			out = append(out, e)
		}
	}
	return out
}
func Since(entries []model.TimelineEntry, since time.Time) []model.TimelineEntry {
	out := []model.TimelineEntry{}
	for _, e := range entries {
		if !e.At.Before(since) {
			out = append(out, e)
		}
	}
	return out
}
func Sort(entries []model.TimelineEntry) []model.TimelineEntry {
	out := append([]model.TimelineEntry(nil), entries...)
	sort.SliceStable(out, func(i, j int) bool { return out[i].At.Before(out[j].At) })
	return out
}
func Search(entries []model.TimelineEntry, token string) []model.TimelineEntry {
	token = strings.ToLower(token)
	out := []model.TimelineEntry{}
	for _, e := range entries {
		if strings.Contains(strings.ToLower(e.Detail), token) {
			out = append(out, e)
		}
	}
	return out
}
