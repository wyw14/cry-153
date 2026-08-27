package model

import "time"

type Clock struct{ now func() time.Time }

func NewClock() Clock { return Clock{now: time.Now} }
func (c Clock) Now() time.Time {
	if c.now == nil {
		return time.Now()
	}
	return c.now()
}
func WindowReady(readings []Reading, threshold float64, span time.Duration, now time.Time) RecoveryDecision {
	if len(readings) == 0 {
		return RecoveryDecision{Reason: "no samples"}
	}
	latest := readings[len(readings)-1]
	if now.Sub(latest.At) > span {
		return RecoveryDecision{Reason: "sample window stale"}
	}
	total := 0.0
	for _, r := range readings {
		total += r.Concentration
	}
	avg := total / float64(len(readings))
	if avg > threshold {
		return RecoveryDecision{Score: avg, Reason: "concentration above threshold"}
	}
	return RecoveryDecision{Ready: true, Score: avg, Reason: "window stable"}
}
