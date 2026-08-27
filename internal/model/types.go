package model

import "time"

type IncidentState string

const (
	StateSuspected   IncidentState = "suspected"
	StateConfirmed   IncidentState = "confirmed"
	StatePropagating IncidentState = "propagating"
	StateContained   IncidentState = "contained"
	StateRecovering  IncidentState = "recovering"
	StateClosed      IncidentState = "closed"
)

type Incident struct {
	ID        string
	State     IncidentState
	CreatedAt time.Time
	UpdatedAt time.Time
	Source    string
	Notes     []string
}
type Reading struct {
	StationID     string
	At            time.Time
	Concentration float64
	Quality       string
}
type SampleBatch struct {
	IncidentID string
	Readings   []Reading
	ReceivedAt time.Time
}
type Intake struct {
	ID        string
	Name      string
	Open      bool
	LockedBy  string
	UpdatedAt time.Time
}
type GateCommand struct {
	ID         string
	IntakeID   string
	IncidentID string
	Close      bool
	IssuedAt   time.Time
	Payload    []byte
}
type Alert struct {
	ID         string
	IncidentID string
	IntakeID   string
	Active     bool
	Message    string
	CreatedAt  time.Time
}
type Operation struct {
	ID         string
	Kind       string
	IncidentID string
	Revision   string
	AcceptedAt time.Time
}
type TimelineEntry struct {
	At     time.Time
	Kind   string
	Detail string
}
type SampleWindow struct {
	StationID string
	Readings  []Reading
}
type RecoveryDecision struct {
	Ready  bool
	Score  float64
	Reason string
}
