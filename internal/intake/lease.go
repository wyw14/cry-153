package intake

import (
	"context"
	"fmt"
	"time"
)

type Lease struct {
	ID    string
	Owner string
	Until time.Time
}

func AcquireLease(ctx context.Context, id, owner string, ttl time.Duration) (Lease, error) {
	if id == "" || owner == "" {
		return Lease{}, fmt.Errorf("lease identity required")
	}
	if ttl <= 0 {
		ttl = time.Minute
	}
	select {
	case <-ctx.Done():
		return Lease{}, ctx.Err()
	default:
		return Lease{ID: id, Owner: owner, Until: time.Now().Add(ttl)}, nil
	}
}
func (l Lease) Valid(now time.Time) bool { return l.ID != "" && l.Owner != "" && now.Before(l.Until) }
func Renew(l Lease, ttl time.Duration) Lease {
	if ttl <= 0 {
		ttl = time.Minute
	}
	l.Until = time.Now().Add(ttl)
	return l
}
