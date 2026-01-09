package poller

import (
	"errors"
	"time"
)

const (
	defaultWorkers  = 3
	defaultInterval = 10 * time.Millisecond
)

// Config defines poller's configuration.
type Config struct {
	Address   string
	WorkerNum int
	Interval  time.Duration
}

// NewConfig creates a new config for poller.
func NewConfig(address string) (*Config, error) {
	if address == "" {
		return nil, errors.New("could not create poller config as the address is empty")
	}
	return &Config{
		Address:   address,
		WorkerNum: defaultWorkers,
		Interval:  defaultInterval,
	}, nil
}
