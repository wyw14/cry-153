package sample

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/wyw14/cry-153/internal/station"
)

// TestCoordinatorOneSourceDisconnectDoesNotPanic reproduces the production
// fault: when one upstream monitoring station disconnects, the shared output
// channel used to be closed, and a still-online station's next send panicked
// with "send on closed channel".
//
// Station A streams far longer than B. B disconnects first; before the fix
// its exit closed the shared output channel and A's next send panicked. With
// the fix, A keeps streaming and the test completes without any goroutine
// panic-crashing the process.
func TestCoordinatorOneSourceDisconnectDoesNotPanic(t *testing.T) {
	reg := station.NewRegistry()
	c := NewCoordinator(reg, 8)
	out := c.Output()

	// Drain the output so online stations can keep sending without blocking.
	drained := make(chan struct{})
	go func() {
		for range out {
		}
		close(drained)
	}()
	defer func() {
		c.Close()
		<-drained
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// A cancellation-based source that stays online until we stop it, so its
	// eventual offline state is unambiguously caused by us — not by a panic.
	aliveCtx, aliveCancel := context.WithCancel(context.Background())
	defer aliveCancel()

	c.StartSource(ctx, station.Source{ID: "B", Values: []float64{10}}) // disconnects first
	// A streams many values so it stays online throughout the observation
	// window after B disconnects (300 values * 1ms = ~300ms > 120ms check).
	c.StartSource(aliveCtx, station.Source{ID: "A", Values: longValues(300), Delay: time.Millisecond})

	// Wait for B to go offline while A is still streaming. Before the fix, A
	// would have panicked here and taken the whole process down.
	deadline := time.After(2 * time.Second)
	for {
		select {
		case <-deadline:
			t.Fatalf("timed out waiting for station B to go offline")
		default:
		}
		if status := reg.Status(); !status["B"] {
			break
		}
		time.Sleep(time.Millisecond)
	}

	// Keep A streaming well past B's disconnect to exercise the send path that
	// used to panic on the now-closed shared channel.
	time.Sleep(60 * time.Millisecond)
	if status := reg.Status(); !status["A"] {
		t.Fatalf("station A went offline after B disconnected (send on closed channel)")
	}

	// Stop A ourselves; it must then go offline cleanly.
	aliveCancel()
	time.Sleep(30 * time.Millisecond)
	if status := reg.Status(); status["A"] {
		t.Fatalf("station A still online after its context was cancelled")
	}
}

// TestCoordinatorCloseDropsInFlightSafely ensures Close does not race with
// active senders: no "send on closed channel" panic, and the output channel
// is closed exactly once so consumers drain and terminate.
func TestCoordinatorCloseDropsInFlightSafely(t *testing.T) {
	reg := station.NewRegistry()
	c := NewCoordinator(reg, 2)
	out := c.Output()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// A source streaming while Close runs concurrently underneath it.
	c.StartSource(ctx, station.Source{ID: "C", Values: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, Delay: 2 * time.Millisecond})

	c.Close()

	// Draining must terminate: the channel is closed exactly once.
	for range out {
	}
}

// TestCoordinatorCloseIsIdempotent ensures repeated Close calls never panic
// (no double close of the shared output channel).
func TestCoordinatorCloseIsIdempotent(t *testing.T) {
	c := NewCoordinator(station.NewRegistry(), 1)
	c.Close()
	c.Close()
}

// TestCoordinatorConcurrentSourcesCountReadings is a stress check: many short
// sources must not panic on each other's disconnect, and all readings must be
// delivered to the consumer.
func TestCoordinatorConcurrentSourcesCountReadings(t *testing.T) {
	const n = 16
	reg := station.NewRegistry()
	c := NewCoordinator(reg, n*4)
	out := c.Output()

	var got atomic.Int64
	done := make(chan struct{})
	go func() {
		for range out {
			got.Add(1)
		}
		close(done)
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	want := int64(0)
	for i := 0; i < n; i++ {
		vals := []float64{1, 2, 3}
		want += int64(len(vals))
		c.StartSource(ctx, station.Source{ID: stationID(i), Values: vals})
	}

	// Wait for every source to go offline.
	deadline := time.After(2 * time.Second)
	for {
		select {
		case <-deadline:
			t.Fatalf("timed out waiting for all sources to disconnect")
		default:
		}
		offline := 0
		for _, v := range reg.Status() {
			if !v {
				offline++
			}
		}
		if offline == n {
			break
		}
		time.Sleep(time.Millisecond)
	}

	c.Close()
	<-done
	if got.Load() != want {
		t.Fatalf("delivered %d readings, want %d", got.Load(), want)
	}
}

func stationID(i int) string {
	return string(rune('A' + i))
}

// longValues returns n monotonically increasing concentrations so a station
// streams long enough to outlive a concurrent short station's disconnect.
func longValues(n int) []float64 {
	out := make([]float64, n)
	for i := range out {
		out[i] = float64(i + 1)
	}
	return out
}
