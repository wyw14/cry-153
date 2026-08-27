package gate

import (
	"context"
	"time"
)

// BoundedBackoff returns an exponentially growing backoff for attempt that is
// always clamped to [0, max]. A negative or overflowed duration must never
// escape this function: once the exponent overflows, time.Duration
// multiplication wraps — sometimes to a negative value, sometimes to zero —
// and time.NewTimer then fires immediately, driving the request rate to
// thousands per second. The observed "negative retry wait time" in the logs
// is exactly that overflow.
func BoundedBackoff(attempt int, max time.Duration) time.Duration {
	if max <= 0 {
		return 0
	}
	if attempt < 0 {
		return 0
	}
	// 1<<shift seconds overflows time.Duration (an int64 of nanoseconds) for
	// shift >= 34. Cap well below that: 30 is ~34 years, far beyond any sane
	// max, and keeps the multiplication from wrapping to a negative or zero
	// duration that would slip past a simple sign check.
	shift := attempt
	if shift > 30 {
		return max
	}
	d := time.Second * time.Duration(1<<shift)
	// Overflow guards: clamp rather than trust the multiplication.
	if d < 0 || d > max {
		return max
	}
	return d
}
func WaitBounded(ctx context.Context, attempt int, max time.Duration) error {
	d := BoundedBackoff(attempt, max)
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}
