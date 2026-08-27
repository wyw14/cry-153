package main

import "time"

type Config struct {
	Address          string
	SnapshotInterval time.Duration
	MaxInFlight      int
}

func DefaultConfig() Config {
	return Config{Address: ":21253", SnapshotInterval: 10 * time.Second, MaxInFlight: 4}
}
