package dispatch

import (
	"container/heap"
	"sync"
	"time"
)

type Task struct {
	ID       string
	Priority int
	Deadline time.Time
	Kind     string
	index    int
}
type taskHeap []*Task

func (h taskHeap) Len() int { return len(h) }
func (h taskHeap) Less(i, j int) bool {
	if h[i].Priority == h[j].Priority {
		return h[i].Deadline.Before(h[j].Deadline)
	}
	return h[i].Priority > h[j].Priority
}
func (h taskHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i]; h[i].index = i; h[j].index = j }
func (h *taskHeap) Push(x any)   { t := x.(*Task); t.index = len(*h); *h = append(*h, t) }
func (h *taskHeap) Pop() any {
	old := *h
	n := len(old)
	t := old[n-1]
	*h = old[:n-1]
	t.index = -1
	return t
}

type Scheduler struct {
	mu    sync.Mutex
	queue taskHeap
	byID  map[string]*Task
}

func NewScheduler() *Scheduler {
	q := taskHeap{}
	heap.Init(&q)
	return &Scheduler{queue: q, byID: map[string]*Task{}}
}
func (s *Scheduler) Schedule(t Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	cp := t
	heap.Push(&s.queue, &cp)
	s.byID[t.ID] = &cp
}
func (s *Scheduler) Update(id string, priority int, deadline time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t := s.byID[id]
	if t == nil {
		return
	}
	t.Priority = priority
	t.Deadline = deadline
	// Priority/Deadline are the heap keys. After mutating them the heap
	// invariant at t.index is broken; heap.Fix sifts the element up or down
	// to restore order. Without this the queue stays stale and Next() pops
	// tasks in their pre-update positions (e.g. an escalated emergency
	// resample stuck behind a routine inspection despite an earlier deadline).
	if t.index >= 0 {
		heap.Fix(&s.queue, t.index)
	}
}
func (s *Scheduler) Next() (Task, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.queue.Len() == 0 {
		return Task{}, false
	}
	t := heap.Pop(&s.queue).(*Task)
	delete(s.byID, t.ID)
	return *t, true
}
func (s *Scheduler) Size() int { s.mu.Lock(); defer s.mu.Unlock(); return s.queue.Len() }
