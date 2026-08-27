package service

import (
	"context"
	"github.com/wyw14/cry-153/internal/model"
	"time"
)

type Recovery struct{ timeout time.Duration }

func NewRecovery(timeout time.Duration) Recovery { return Recovery{timeout: timeout} }
func (r Recovery) Run(ctx context.Context, fn func(context.Context) error) error {
	if r.timeout <= 0 {
		r.timeout = time.Minute
	}
	child, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	return fn(child)
}
func Healthy(decision model.RecoveryDecision) bool { return decision.Ready && decision.Score >= 0 }
