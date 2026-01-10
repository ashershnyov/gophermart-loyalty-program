package poller

import (
	"errors"
	"time"
)

const (
	defaultWorkers  = 3
	defaultInterval = 1 * time.Second
)

// Config defines poller's configuration.
type Config struct {
	Address   string
	WorkerNum int
	Interval  time.Duration
}

// Option is a config option setter.
type Option func(*Config)

// NewConfig creates a new config for poller.
func NewConfig(opts ...Option) (*Config, error) {
	c := &Config{
		WorkerNum: defaultWorkers,
		Interval:  defaultInterval,
	}

	for _, o := range opts {
		o(c)
	}

	if c.Address == "" {
		return nil, errors.New("could not create poller config as the address is empty")
	}

	return c, nil
}

// SetAddress sets the address for the accrual service to make requests to.
func SetAddress(addr string) Option {
	return func(c *Config) {
		c.Address = addr
	}
}

// SetWorkerNum sets the amount of poller's workers.
func SetWorkerNum(num int) Option {
	return func(c *Config) {
		c.WorkerNum = num
	}
}

// SetInterval sets the polling interval.
func SetInterval(interval time.Duration) Option {
	return func(c *Config) {
		c.Interval = interval
	}
}
