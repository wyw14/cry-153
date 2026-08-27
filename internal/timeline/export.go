package timeline

import (
	"encoding/json"
	"github.com/wyw14/cry-153/internal/model"
	"sort"
)

func Encode(entries []model.TimelineEntry) ([]byte, error) {
	cp := append([]model.TimelineEntry(nil), entries...)
	sort.SliceStable(cp, func(i, j int) bool { return cp[i].At.Before(cp[j].At) })
	return json.Marshal(cp)
}
func Decode(data []byte) ([]model.TimelineEntry, error) {
	var out []model.TimelineEntry
	err := json.Unmarshal(data, &out)
	return out, err
}
func CountByKind(entries []model.TimelineEntry) map[string]int {
	out := map[string]int{}
	for _, e := range entries {
		out[e.Kind]++
	}
	return out
}
func Has(entries []model.TimelineEntry, kind string) bool {
	for _, e := range entries {
		if e.Kind == kind {
			return true
		}
	}
	return false
}
