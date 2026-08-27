package journal

import (
	"fmt"
	"github.com/wyw14/cry-153/internal/model"
	"sync"
	"time"
)

type Logger struct {
	mu      sync.Mutex
	records []string
}

func NewLogger() *Logger { return &Logger{} }
func (l *Logger) Record(kind, detail string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.records = append(l.records, fmt.Sprintf("%s %s %s", time.Now().UTC().Format(time.RFC3339Nano), kind, detail))
}
func (l *Logger) Records() []string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return append([]string(nil), l.records...)
}
func Timeline(kind, detail string) model.TimelineEntry {
	return model.TimelineEntry{At: time.Now(), Kind: kind, Detail: detail}
}
