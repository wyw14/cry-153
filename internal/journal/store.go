package journal

import (
	"encoding/json"
	"github.com/wyw14/cry-153/internal/model"
	"os"
	"sync"
)

type Store struct {
	mu      sync.Mutex
	path    string
	entries []model.TimelineEntry
}

func NewStore(path string) *Store { return &Store{path: path} }
func (s *Store) Append(entry model.TimelineEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = append(s.entries, entry)
	return s.flush()
}
func (s *Store) Entries() []model.TimelineEntry {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]model.TimelineEntry, len(s.entries))
	copy(out, s.entries)
	return out
}
func (s *Store) flush() error {
	if s.path == "" {
		return nil
	}
	data, err := json.Marshal(s.entries)
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err = os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}
