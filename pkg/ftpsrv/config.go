package ftpsrv

import (
	"time"
)

// Option represents an FTP configuration option.
type Option func(Config) Config

// WithHost is the [NewConfig] option setting the host to listen on.
func WithHost(host string) Option {
	return func(cfg Config) Config {
		cfg.host = host
		return cfg
	}
}

// WithPort is the [NewConfig] option setting the port to listen on.
func WithPort(port int) Option {
	return func(cfg Config) Config {
		cfg.port = port
		return cfg
	}
}

// WithClock is the [NewConfig] option setting clock to use.
func WithClock(clk func() time.Time) Option {
	return func(cfg Config) Config {
		cfg.clock = clk
		return cfg
	}
}

// WithTLS is the [NewConfig] option setting TLS certificate and key.
func WithTLS(cert, key []byte) Option {
	return func(cfg Config) Config {
		cfg.supportTLS = true
		cfg.cert = cert
		cfg.key = key
		return cfg
	}
}

// WithImpTLS is the [NewConfig] option turning on implicit TLS.
func WithImpTLS(cfg Config) Config {
	cfg.implicitTLS = true
	return cfg
}

// WithReadyMsg is the [NewConfig] option setting custom server ready response.
func WithReadyMsg(rsp Response) Option {
	return func(cfg Config) Config {
		cfg.svrReadyMsg = rsp
		return cfg
	}
}

func WithReadTimeout(timeout time.Duration) Option {
	return func(cfg Config) Config {
		cfg.readTO = timeout
		return cfg
	}
}

func WithWriteTimeout(timeout time.Duration) Option {
	return func(cfg Config) Config {
		cfg.writeTO = timeout
		return cfg
	}
}

// Config represents FTP server configuration.
type Config struct {
	host        string           // FTP listening address.
	port        int              // FTP listening port.
	readTO      time.Duration    // Connection read timeout.
	writeTO     time.Duration    // Connection write timeout.
	clock       func() time.Time // Current time and timezone.
	supportTLS  bool             // Support TLS connections (AUTH TLS command).
	implicitTLS bool             // Use implicit TLS.
	svrReadyMsg Response         // Server ready message.

	cert []byte // Certificate PEM block.
	key  []byte // Certificate key PEM block.
}

// NewConfig returns default FTP server configuration options.
func NewConfig(opts ...Option) Config {
	cfg := Config{
		host:        "127.0.0.1",
		port:        21,
		readTO:      150 * time.Millisecond,
		writeTO:     50 * time.Millisecond,
		clock:       func() time.Time { return time.Now().UTC() },
		svrReadyMsg: ServerReady,
	}
	for _, opt := range opts {
		cfg = opt(cfg)
	}
	return cfg
}
