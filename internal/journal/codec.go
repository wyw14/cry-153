package journal

import (
	"encoding/json"
	"github.com/wyw14/cry-153/internal/model"
)

func Encode(entries []model.TimelineEntry) ([]byte, error) { return json.Marshal(entries) }
func Decode(data []byte) ([]model.TimelineEntry, error) {
	var out []model.TimelineEntry
	err := json.Unmarshal(data, &out)
	return out, err
}
func Clone(entries []model.TimelineEntry) []model.TimelineEntry {
	out := make([]model.TimelineEntry, len(entries))
	copy(out, entries)
	return out
}
func AppendMany(s *Store, entries []model.TimelineEntry) error {
	for _, e := range entries {
		if err := s.Append(e); err != nil {
			return err
		}
	}
	return nil
}
