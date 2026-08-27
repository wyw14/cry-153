package model

import (
	"encoding/json"
	"fmt"
	"time"
)

func MarshalIncident(item Incident) ([]byte, error) { return json.Marshal(item) }
func FormatReading(r Reading) string {
	return fmt.Sprintf("%s %.3f %s", r.StationID, r.Concentration, r.At.UTC().Format(time.RFC3339))
}
func CloneReadings(in []Reading) []Reading {
	out := make([]Reading, len(in))
	copy(out, in)
	return out
}
func MergeReadings(a, b []Reading) []Reading {
	out := make([]Reading, 0, len(a)+len(b))
	out = append(out, a...)
	out = append(out, b...)
	return out
}
