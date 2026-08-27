package model

import "time"

type Deadline struct {
	At    time.Time
	Grace time.Duration
}

func NewDeadline(at time.Time, grace time.Duration) Deadline { return Deadline{At: at, Grace: grace} }
func (d Deadline) Expired(now time.Time) bool                { return now.After(d.At.Add(d.Grace)) }
func (d Deadline) Remaining(now time.Time) time.Duration {
	left := d.At.Add(d.Grace).Sub(now)
	if left < 0 {
		return 0
	}
	return left
}
func Later(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}
func Earlier(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}
func RoundToSecond(t time.Time) time.Time { return t.Truncate(time.Second) }
