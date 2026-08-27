package dispatch

import (
	"sort"
	"time"
)

func Order(tasks []Task) []Task {
	out := append([]Task(nil), tasks...)
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Priority == out[j].Priority {
			return out[i].Deadline.Before(out[j].Deadline)
		}
		return out[i].Priority > out[j].Priority
	})
	return out
}
func Due(tasks []Task, now time.Time) []Task {
	out := []Task{}
	for _, t := range tasks {
		if !t.Deadline.After(now) {
			out = append(out, t)
		}
	}
	return out
}
func NextDeadline(tasks []Task) (time.Time, bool) {
	if len(tasks) == 0 {
		return time.Time{}, false
	}
	d := tasks[0].Deadline
	for _, t := range tasks[1:] {
		if t.Deadline.Before(d) {
			d = t.Deadline
		}
	}
	return d, true
}
