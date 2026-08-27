package model

import (
	"fmt"
	"github.com/google/uuid"
)

func NewID(prefix string) string { return fmt.Sprintf("%s-%s", prefix, uuid.NewString()) }
func NewRevision() string        { return uuid.NewString() }
