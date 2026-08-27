package station

import (
	"context"
	"github.com/wyw14/cry-153/internal/model"
	"time"
)

type Source struct {
	ID     string
	Values []float64
	Delay  time.Duration
}

func (s Source) Stream(ctx context.Context, emit func(model.Reading)) error {
	for _, value := range s.Values {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		emit(FreshReading(s.ID, value))
		if s.Delay > 0 {
			timer := time.NewTimer(s.Delay)
			select {
			case <-ctx.Done():
				timer.Stop()
				return ctx.Err()
			case <-timer.C:
			}
		}
	}
	return nil
}
