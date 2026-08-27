package intake

import (
	"github.com/wyw14/cry-153/internal/model"
	"sort"
	"sync"
)

type Locks struct {
	mu      sync.Mutex
	items   map[string]*sync.Mutex
	intakes map[string]model.Intake
}

func NewLocks(ids []string) *Locks {
	l := &Locks{items: map[string]*sync.Mutex{}, intakes: map[string]model.Intake{}}
	for _, id := range ids {
		l.items[id] = &sync.Mutex{}
		l.intakes[id] = model.Intake{ID: id, Name: id, Open: true}
	}
	return l
}
func (l *Locks) Ordered(ids []string) []string {
	out := append([]string(nil), ids...)
	sort.Strings(out)
	return out
}
func (l *Locks) Lock(id string)   { l.items[id].Lock() }
func (l *Locks) Unlock(id string) { l.items[id].Unlock() }
func (l *Locks) Set(id string, open bool) {
	l.mu.Lock()
	x := l.intakes[id]
	x.Open = open
	l.intakes[id] = x
	l.mu.Unlock()
}
func (l *Locks) Snapshot() map[string]model.Intake {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := map[string]model.Intake{}
	for k, v := range l.intakes {
		out[k] = v
	}
	return out
}
