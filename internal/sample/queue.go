package sample

import (
	"context"
	"github.com/wyw14/cry-153/internal/model"
	"sync"
)

type Queue struct {
	mu     sync.Mutex
	items  []model.Reading
	closed bool
}

func NewQueue() *Queue { return &Queue{items: []model.Reading{}} }
func (q *Queue) Push(r model.Reading) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return false
	}
	q.items = append(q.items, r)
	return true
}
func (q *Queue) Close() { q.mu.Lock(); q.closed = true; q.mu.Unlock() }
func (q *Queue) Drain(ctx context.Context) []model.Reading {
	q.mu.Lock()
	defer q.mu.Unlock()
	out := model.CloneReadings(q.items)
	q.items = q.items[:0]
	select {
	case <-ctx.Done():
		return nil
	default:
		return out
	}
}
func (q *Queue) Len() int { q.mu.Lock(); defer q.mu.Unlock(); return len(q.items) }
