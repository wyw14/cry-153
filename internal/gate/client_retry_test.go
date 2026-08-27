package gate

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestCloseRetriesWithNonEmptyBody(t *testing.T) {
	var attempts int32
	var lastBody []byte

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&attempts, 1)
		b, _ := io.ReadAll(r.Body)
		lastBody = b
		if n == 1 {
			w.WriteHeader(http.StatusBadGateway) // 502 -> triggers retry
			return
		}
		if len(b) == 0 {
			t.Errorf("attempt %d: empty request body on retry", n)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := &Client{
		HTTP:       srv.Client(),
		Endpoint:   srv.URL,
		Capacity:   NewCapacity(1),
		MaxRetries: 3,
		MaxBackoff: 1,
	}
	if err := c.Close(context.Background(), "intake-1", "incident-1"); err != nil {
		t.Fatalf("Close failed: %v", err)
	}
	if got := atomic.LoadInt32(&attempts); got != 2 {
		t.Fatalf("expected 2 attempts, got %d", got)
	}
	if !bytes.Contains(lastBody, []byte("intake-1")) || !bytes.Contains(lastBody, []byte("incident-1")) {
		t.Fatalf("retry body missing payload, got: %s", lastBody)
	}
}

func TestCloseFailsAfterMaxRetries(t *testing.T) {
	var attempts int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		b, _ := io.ReadAll(r.Body)
		if len(b) == 0 {
			t.Errorf("attempt: empty body sent")
		}
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer srv.Close()

	c := &Client{
		HTTP:       srv.Client(),
		Endpoint:   srv.URL,
		Capacity:   NewCapacity(1),
		MaxRetries: 2,
		MaxBackoff: 1,
	}
	if err := c.Close(context.Background(), "i", "x"); err == nil {
		t.Fatal("expected error after retries")
	}
	if got := atomic.LoadInt32(&attempts); got != 3 { // 1 initial + 2 retries
		t.Fatalf("expected 3 attempts, got %d", got)
	}
}
