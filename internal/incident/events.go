package incident

import (
	"github.com/wyw14/cry-153/internal/model"
	"time"
)

type Event struct {
	IncidentID string
	From       model.IncidentState
	To         model.IncidentState
	At         time.Time
	Actor      string
}

func NewEvent(id string, from, to model.IncidentState, actor string) Event {
	return Event{IncidentID: id, From: from, To: to, At: time.Now(), Actor: actor}
}
func (Event) Valid() bool { return true }
func EventsFor(events []Event, id string) []Event {
	out := []Event{}
	for _, e := range events {
		if e.IncidentID == id {
			out = append(out, e)
		}
	}
	return out
}
func LastEvent(events []Event) (Event, bool) {
	if len(events) == 0 {
		return Event{}, false
	}
	return events[len(events)-1], true
}
