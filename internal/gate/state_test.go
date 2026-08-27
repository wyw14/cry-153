package gate

import (
	"context"
	"sync"
	"testing"
	"time"
)

// TestCapacityCancelledWaitsDoNotInflateInFlight reproduces the concurrency
// fault: repeatedly cancelling Acquire requests that are blocked WAITING for a
// slot must not inflate the inFlight counter, because those callers never
// actually acquired a token (and never call Release).
//
// Before the fix, each cancelled wait bumped inFlight by 1 even though no
// token was consumed. With a limit of 4 and 3 cancelled waits, inFlight climbed
// from 4 to 7 while the token channel stayed correctly sized at 4 -- leaving
// the gauge permanently desynchronised from the real semaphore.
func TestCapacityCancelledWaitsDoNotInflateInFlight(t *testing.T) {
	const limit = 4
	cap := NewCapacity(limit)

	// Saturate every token so subsequent Acquires block while waiting.
	type holder struct {
		release context.CancelFunc
	}
	holders := make([]holder, 0, limit)
	for i := 0; i < limit; i++ {
		ctx, cancel := context.WithCancel(context.Background())
		if err := cap.Acquire(ctx); err != nil {
			t.Fatalf("initial acquire %d: %v", i, err)
		}
		holders = append(holders, holder{release: cancel})
	}
	if got := cap.InFlight(); got != limit {
		t.Fatalf("after saturation inFlight=%d want %d", got, limit)
	}

	// Repeatedly cancel Acquire requests that are blocked waiting for a slot.
	const cancelled = 3
	var wg sync.WaitGroup
	for i := 0; i < cancelled; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithCancel(context.Background())
			go cancel() // cancel almost immediately while waiting
			_ = cap.Acquire(ctx) // expected to return ctx.Err()
		}()
	}
	wg.Wait()

	// Give any lingering goroutines a moment to settle.
	time.Sleep(20 * time.Millisecond)

	// The cancelled waits never consumed a token, so inFlight must still equal
	// the number of live token holders -- NOT limit+cancelled.
	if got := cap.InFlight(); got != limit {
		t.Fatalf("cancelled waits inflated inFlight: got %d want %d (limit)", got, limit)
	}

	// Release all real holders; inFlight must drain to 0.
	for _, h := range holders {
		cap.Release()
		h.release()
	}
	if got := cap.InFlight(); got != 0 {
		t.Fatalf("after release inFlight=%d want 0", got)
	}
}

// TestCapacityRealConcurrencyCappedAfterCancellations verifies the downstream
// effect that protects the gateway: after repeatedly cancelling blocked waits,
// the hard semaphore still admits at most `limit` concurrent acquirers -- it
// never lets the live in-flight count exceed the configured limit. Before the
// fix the inflated gauge could be read as spare capacity, admitting more gate
// close requests in parallel than the gateway could sustain.
func TestCapacityRealConcurrencyCappedAfterCancellations(t *testing.T) {
	const limit = 4
	cap := NewCapacity(limit)

	// Hold all tokens.
	holders := make([]context.CancelFunc, 0, limit)
	for i := 0; i < limit; i++ {
		ctx, cancel := context.WithCancel(context.Background())
		if err := cap.Acquire(ctx); err != nil {
			t.Fatalf("acquire %d: %v", i, err)
		}
		holders = append(holders, cancel)
	}

	// Spam cancelled waits the way the fault report describes.
	var wg sync.WaitGroup
	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithCancel(context.Background())
			go cancel()
			_ = cap.Acquire(ctx)
		}()
	}
	wg.Wait()
	time.Sleep(10 * time.Millisecond)

	if got := cap.InFlight(); got > limit {
		t.Fatalf("inFlight %d exceeds limit %d after cancellations", got, limit)
	}

	// New waits must still block while all tokens are held: the spam must not
	// have punched a hole in the semaphore.
	probeCtx, probeCancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer probeCancel()
	if err := cap.Acquire(probeCtx); err == nil {
		t.Fatalf("probe acquired a token while limit saturated; semaphore leaked")
	}

	// Free one slot; exactly one new acquirer may now proceed.
	holders[0]()
	cap.Release()

	admitCtx, admitCancel := context.WithTimeout(context.Background(), time.Second)
	defer admitCancel()
	if err := cap.Acquire(admitCtx); err != nil {
		t.Fatalf("acquire after free: %v", err)
	}
	cap.Release()
	for _, c := range holders[1:] {
		c()
		cap.Release()
	}
	if got := cap.InFlight(); got != 0 {
		t.Fatalf("final inFlight=%d want 0", got)
	}
}
