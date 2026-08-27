package alert

import (
	"context"
	"github.com/wyw14/cry-153/internal/model"
	"sync"
)

type Queue struct {
	mu     sync.Mutex
	items  []Message
	closed bool
}

func NewQueue() *Queue { return &Queue{items: []Message{}} }
func (q *Queue) Push(m Message) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return false
	}
	q.items = append(q.items, m)
	return true
}
func (q *Queue) Close() { q.mu.Lock(); q.closed = true; q.mu.Unlock() }
func (q *Queue) Drain(ctx context.Context) []Message {
	q.mu.Lock()
	defer q.mu.Unlock()
	select {
	case <-ctx.Done():
		return nil
	default:
		out := append([]Message(nil), q.items...)
		q.items = q.items[:0]
		return out
	}
}
func (q *Queue) Len() int { q.mu.Lock(); defer q.mu.Unlock(); return len(q.items) }
func ToAlert(m Message) model.Alert {
	return model.Alert{ID: model.NewID("alert"), IncidentID: m.IncidentID, IntakeID: m.IntakeID, Active: m.Active, Message: m.Text}
}
