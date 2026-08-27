package plume

import (
	"errors"
	"fmt"
	"strings"
)

type Classification string

const (
	ClassTemporary Classification = "temporary"
	ClassPermanent Classification = "permanent"
)

func Classify(err error) Classification {
	if errors.Is(err, ErrTemporaryUpstream) {
		return ClassTemporary
	}
	return ClassPermanent
}
func Wrap(err error, stage string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", stage, err)
}
func IsRetryable(err error) bool { return Classify(err) == ClassTemporary }
func RetryableForLocale(m *Model, err error) bool { return strings.Contains(m.Display(err), "temporary") }
