package service

import (
	"github.com/wyw14/cry-153/internal/model"
	"sync"
)

type Operations struct {
	mu    sync.RWMutex
	items map[string]model.Operation
}

func NewOperations() *Operations               { return &Operations{items: map[string]model.Operation{}} }
func (o *Operations) Put(item model.Operation) { o.mu.Lock(); o.items[item.ID] = item; o.mu.Unlock() }
func (o *Operations) Get(id string) (model.Operation, bool) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	item, ok := o.items[id]
	return item, ok
}
func (o *Operations) List() []model.Operation {
	o.mu.RLock()
	defer o.mu.RUnlock()
	out := make([]model.Operation, 0, len(o.items))
	for _, item := range o.items {
		out = append(out, item)
	}
	return out
}
