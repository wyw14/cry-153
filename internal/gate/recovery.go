package gate

import (
	"context"
	"time"
)

func BoundedBackoff(attempt int, max time.Duration) time.Duration {
	if attempt < 0 {
		return 0
	}
	if attempt > 30 {
		return max
	}
	d := time.Second * time.Duration(uint64(1)<<uint(attempt))
	if d > max {
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
