package ftpsrv

import (
	"maps"
	"sort"
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
	// TODO(rz): document this.
	// TODO(rz): test this.
	return func(cfg Config) Config {
		cfg.readTO = timeout
		return cfg
	}
}

func WithWriteTimeout(timeout time.Duration) Option {
	// TODO(rz): document this.
	// TODO(rz): test this.
	return func(cfg Config) Config {
		cfg.writeTO = timeout
		return cfg
	}
}

// WithoutFeature is the [NewConfig] option disabling support for FTP command.
func WithoutFeature(feat string) Option {
	return func(cfg Config) Config { return cfg.DisableFeature(feat) }
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

	// List of available FTP commands.
	features map[string]Command
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
		features:    make(map[string]Command, 20),
	}
	// Get all implemented commands.
	for _, cmd := range commands {
		cfg.features[cmd.Name()] = cmd
	}
	for _, opt := range opts {
		cfg = opt(cfg)
	}
	return cfg
}

// Features returns a sorted slice of enabled features.
func (cfg Config) Features() []string {
	var feats []string
	for cmd := range cfg.features {
		feats = append(feats, cmd)
	}
	sort.Strings(feats)
	return feats
}

// HasFeature returns true when the feature is enabled.
func (cfg Config) HasFeature(feat string) bool {
	_, ok := cfg.features[feat]
	return ok
}

// TODO(rz):
// // EnableFeature enables given FTP command(s).
// func (cfg Config) EnableFeature(commands ...string) Config {
// 	cfg.features = maps.Clone(cfg.features)
// 	for _, cmd := range commands {
// 		cfg.features[cmd] = struct{}{}
// 	}
// 	return cfg
// }

// DisableFeature disables the given FTP command.
func (cfg Config) DisableFeature(name string) Config {
	cfg.features = maps.Clone(cfg.features)
	delete(cfg.features, name)
	return cfg
}

// DisableAllFeatures disables all currently enabled features. The original
// instance of Config is not changed in any way.
func (cfg Config) DisableAllFeatures() Config {
	cfg.features = maps.Clone(cfg.features)
	for name := range cfg.features {
		delete(cfg.features, name)
	}
	return cfg
}
