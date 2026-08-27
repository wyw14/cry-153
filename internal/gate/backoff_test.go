package gate

import (
	"context"
	"testing"
	"time"
)

func TestBoundedBackoffNeverNegative(t *testing.T) {
	max := 30 * time.Second
	// Every attempt, including values that previously overflowed the signed
	// int shift, must stay within [0, max]. The old implementation returned a
	// negative duration for attempt >= 63, which made the retry timer fire
	// immediately and produced the thousands-per-second request storm.
	for attempt := -1; attempt < 200; attempt++ {
		d := BoundedBackoff(attempt, max)
		if d < 0 {
			t.Fatalf("attempt %d: backoff %v is negative", attempt, d)
		}
		if d > max {
			t.Fatalf("attempt %d: backoff %v exceeds max %v", attempt, d, max)
		}
	}
}

func TestBoundedBackoffClampsToMax(t *testing.T) {
	max := 30 * time.Second
	// Large attempts must saturate at max instead of overflowing.
	if d := BoundedBackoff(100, max); d != max {
		t.Fatalf("attempt 100: want %v, got %v", max, d)
	}
	if d := BoundedBackoff(63, max); d != max {
		t.Fatalf("attempt 63 (overflow boundary): want %v, got %v", max, d)
	}
}

func TestBoundedBackoffRespectsSmallMax(t *testing.T) {
	// A tiny max must clamp even small attempts. The old code ignored max
	// entirely and returned 1s, 2s, 4s, ... regardless of the bound.
	max := 500 * time.Millisecond
	if d := BoundedBackoff(0, max); d != max {
		t.Fatalf("attempt 0 with tiny max: want %v, got %v", max, d)
	}
}

func TestBoundedBackoffZeroOrNegativeMax(t *testing.T) {
	if d := BoundedBackoff(5, 0); d != 0 {
		t.Fatalf("max 0: want 0, got %v", d)
	}
	if d := BoundedBackoff(5, -time.Second); d != 0 {
		t.Fatalf("negative max: want 0, got %v", d)
	}
}

func TestTotalBackoffBoundedByMax(t *testing.T) {
	max := 30 * time.Second
	// Sum across many retries must never exceed max, even though each raw
	// exponential term would. Previously a negative term slipped past the
	// `max-total < d` cap check and corrupted the running total.
	if total := TotalBackoff(100, max); total > max {
		t.Fatalf("TotalBackoff(100, max) = %v, exceeds max %v", total, max)
	}
	if total := TotalBackoff(100, max); total < 0 {
		t.Fatalf("TotalBackoff(100, max) = %v, negative", total)
	}
}

func TestWaitBackoffDoesNotFireImmediatelyForLargeAttempt(t *testing.T) {
	// The close-loop backoff must honor the bound for overflowing attempts
	// rather than returning instantly (which caused the rate spike).
	max := 50 * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	start := time.Now()
	if err := WaitBackoff(ctx, 100, max); err != nil {
		t.Fatalf("WaitBackoff: %v", err)
	}
	elapsed := time.Since(start)
	if elapsed < max/2 {
		t.Fatalf("WaitBackoff fired after %v, expected to wait near %v", elapsed, max)
	}
}
