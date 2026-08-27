package intake

import (
	"github.com/wyw14/cry-153/internal/model"
	"time"
)

type Audit struct {
	At         time.Time
	IncidentID string
	Intakes    []string
	Accepted   bool
}

func NewAudit(id string, intakes []string, accepted bool) Audit {
	return Audit{At: time.Now(), IncidentID: id, Intakes: Normalize(intakes), Accepted: accepted}
}
func (a Audit) Operation() model.Operation {
	return model.Operation{ID: model.NewID("op"), Kind: "intake-close", IncidentID: a.IncidentID, AcceptedAt: a.At}
}
func MergeAudits(a, b []Audit) []Audit {
	out := append([]Audit(nil), a...)
	out = append(out, b...)
	return out
}
func Recent(audits []Audit, since time.Time) []Audit {
	out := []Audit{}
	for _, a := range audits {
		if !a.At.Before(since) {
			out = append(out, a)
		}
	}
	return out
}
