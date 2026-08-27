package service

import (
	"net/http"
	"sync/atomic"
)

type HealthState struct {
	ready    atomic.Bool
	requests atomic.Uint64
}

func (h *HealthState) MarkReady() { h.ready.Store(true) }
func (h *HealthState) Observe()   { h.requests.Add(1) }
func (h *HealthState) Handler(w http.ResponseWriter, r *http.Request) {
	h.Observe()
	if !h.ready.Load() {
		http.Error(w, "starting", 503)
		return
	}
	w.WriteHeader(200)
	_, _ = w.Write([]byte("ok"))
}
