package model

import (
	"fmt"
	"strings"
	"time"
)

func ValidateIncident(item Incident) error {
	if strings.TrimSpace(item.ID) == "" {
		return fmt.Errorf("incident id required")
	}
	if item.CreatedAt.IsZero() {
		return fmt.Errorf("created time required")
	}
	return nil
}
func ValidateReading(r Reading) error {
	if strings.TrimSpace(r.StationID) == "" {
		return fmt.Errorf("station id required")
	}
	if r.At.IsZero() {
		return fmt.Errorf("reading time required")
	}
	if r.Concentration < 0 {
		return fmt.Errorf("negative concentration")
	}
	return nil
}
func ReadingAge(r Reading, now time.Time) time.Duration {
	if now.Before(r.At) {
		return 0
	}
	return now.Sub(r.At)
}
func Clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
func StateTerminal(s IncidentState) bool { return s == StateClosed }
