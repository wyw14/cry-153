package model

import "time"

func NewOperation(kind, incident string) Operation {
	return Operation{ID: NewID("operation"), Kind: kind, IncidentID: incident, Revision: NewRevision(), AcceptedAt: time.Now()}
}
func (o Operation) Valid() bool { return o.ID != "" && o.Kind != "" && o.IncidentID != "" }
func (o Operation) Age(now time.Time) time.Duration {
	if now.Before(o.AcceptedAt) {
		return 0
	}
	return now.Sub(o.AcceptedAt)
}
func LatestOperation(items []Operation) (Operation, bool) {
	if len(items) == 0 {
		return Operation{}, false
	}
	out := items[0]
	for _, item := range items[1:] {
		if item.AcceptedAt.After(out.AcceptedAt) {
			out = item
		}
	}
	return out, true
}
