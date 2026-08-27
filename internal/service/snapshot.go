package service

import (
	"encoding/json"
	"github.com/wyw14/cry-153/internal/model"
	"os"
	"path/filepath"
)

type Snapshot struct {
	Incidents []model.Incident
	Written   int64
}

func SaveSnapshot(path string, s Snapshot) error {
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
func LoadSnapshot(path string) (Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Snapshot{}, err
	}
	var s Snapshot
	err = json.Unmarshal(data, &s)
	return s, err
}
func EmptySnapshot() Snapshot { return Snapshot{Incidents: []model.Incident{}} }
