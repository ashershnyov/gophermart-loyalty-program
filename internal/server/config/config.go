package config

const (
	// defaultAddress specifies the address used unless overridden by starting params.
	defaultAddress = ":8080"
	// defaultMaxRetries sets the default amount of retries upon send errors.
	defaultMaxRetries = 3
	// defaultDBAddress specifies the address used to connect to the database
	// unless overridden by starting params.
	defaultDBAddress = ""
	// defaultAccrualAddress specifies the address used to connect to the acurral
	// service unless overridden by starting params.
	defaultAccrualAddress = ""
)

// Config specifies the server's configuration.
type Config struct {
	Address        string
	MaxRetries     int
	DBAddress      string
	AccrualAddress string
}

// Option is a config option setter.
type Option func(*Config)

// New returns a newly created config.
func New(opts ...Option) *Config {
	c := &Config{
		Address:        defaultAddress,
		MaxRetries:     defaultMaxRetries,
		DBAddress:      defaultDBAddress,
		AccrualAddress: defaultAccrualAddress,
	}

	for _, o := range opts {
		o(c)
	}
	return c
}

// SetAddress sets the address for the server to run on.
func SetAddress(addr string) Option {
	return func(c *Config) {
		c.Address = addr
	}
}

// SetDBAddress specifies the address to connect to DB on.
func SetDBAddress(addr string) Option {
	return func(c *Config) {
		c.DBAddress = addr
	}
}

// SetAccrualAddress specifies the address to connect to the accrual service on.
func SetAccrualAddress(address string) Option {
	return func(c *Config) {
		c.AccrualAddress = address
	}
}
