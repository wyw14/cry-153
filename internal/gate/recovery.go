package gate

import (
	"context"
	"time"
)

func BoundedBackoff(attempt int, max time.Duration) time.Duration {
	if attempt < 0 {
		return 0
	}
	d := time.Second * time.Duration(1<<attempt)
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
