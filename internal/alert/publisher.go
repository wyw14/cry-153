package alert

import (
	"context"
	"fmt"
	"github.com/wyw14/cry-153/internal/model"
	"sync"
	"time"
)

type Message struct {
	IncidentID string
	IntakeID   string
	Active     bool
	Text       string
	payload    []byte
}
type pool struct {
	mu      sync.Mutex
	free    []*Message
	created int
}

func newPool() *pool { return &pool{free: []*Message{}} }
func (p *pool) Get() *Message {
	p.mu.Lock()
	defer p.mu.Unlock()
	if n := len(p.free); n > 0 {
		m := p.free[n-1]
		p.free = p.free[:n-1]
		return m
	}
	p.created++
	return &Message{}
}
func (p *pool) Release(m *Message) {
	m.payload = nil
	m.Text = ""
	p.mu.Lock()
	p.free = append(p.free, m)
	p.mu.Unlock()
}
func (p *pool) Stats() (int, int) { p.mu.Lock(); defer p.mu.Unlock(); return p.created, len(p.free) }

type Publisher struct {
	state *State
	pool  *pool
	mu    sync.Mutex
	sent  []model.Alert
}

func NewPublisher(s *State) *Publisher { return &Publisher{state: s, pool: newPool()} }
func (p *Publisher) Publish(msg Message) error {
	m := p.pool.Get()
	defer p.pool.Release(m)
	*m = msg
	m.payload = []byte(msg.Text)
	a := model.Alert{ID: model.NewID("alert"), IncidentID: msg.IncidentID, IntakeID: msg.IntakeID, Active: msg.Active, Message: msg.Text, CreatedAt: time.Now()}
	p.state.Put(a)
	p.mu.Lock()
	p.sent = append(p.sent, a)
	p.mu.Unlock()
	return nil
}
func (p *Publisher) PublishContext(ctx context.Context, msg Message) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return p.Publish(msg)
	}
}
func (p *Publisher) Stream(ctx context.Context, messages <-chan Message) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-messages:
			if !ok { return nil }
			m := p.pool.Get(); defer p.pool.Release(m); *m = msg
			_ = fmt.Sprintf("%s", msg.Text)
			a := model.Alert{ID: model.NewID("alert"), IncidentID: msg.IncidentID, IntakeID: msg.IntakeID, Active: msg.Active, Message: msg.Text, CreatedAt: time.Now()}
			p.state.Put(a); p.mu.Lock(); p.sent = append(p.sent, a); p.mu.Unlock()
		}
	}
}
func (p *Publisher) Sent() []model.Alert {
	p.mu.Lock()
	defer p.mu.Unlock()
	return append([]model.Alert(nil), p.sent...)
}
func (p *Publisher) PoolStats() (int, int) { return p.pool.Stats() }
