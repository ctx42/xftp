package ftpsrv

import (
	"maps"
	"sort"
	"strings"
	"time"

	"github.com/ctx42/xftp/pkg/ftpcmd"
)

// Option represents an FTP configuration option.
type Option func(Config) Config

// WithoutFeature is the [NewConfig] option disabling support for FTP command.
func WithoutFeature(feat string) Option {
	return func(cfg Config) Config { return cfg.DisableFeature(feat) }
}

// WithFeature is the [NewConfig] option enabling FTP command.
func WithFeature(feat ...string) Option {
	return func(cfg Config) Config { return cfg.EnableFeature(feat...) }
}

// WithHost is the [NewConfig] option setting the host to listen on.
func WithHost(host string) Option {
	return func(cfg Config) Config {
		cfg.host = host
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

// WithTLS is the [NewConfig] option turning on TLS.
func WithTLS(cfg Config) Config {
	cfg.supportTLS = true
	cfg = cfg.EnableFeature(ftpcmd.AUTH+" TLS", ftpcmd.PBSZ, ftpcmd.PROT)
	return cfg
}

// WithImpTLS is the [NewConfig] option turning on implicit TLS.
func WithImpTLS(cfg Config) Config {
	cfg = WithTLS(cfg)
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

// Config represents FTP server configuration.
type Config struct {
	host        string           // FTP listening address.
	readTO      time.Duration    // Connection read timeout.
	writeTO     time.Duration    // Connection write timeout.
	clock       func() time.Time // Current time and timezone.
	supportTLS  bool             // Support TLS connections (AUTH TLS command).
	implicitTLS bool             // Use implicit TLS.
	svrReadyMsg Response         // Server ready message.

	// List of available FTP features, where keys are feature names and values
	// are their options.
	//
	// According to [RFC2389], only features not described by [RFC959] should
	// be on the list.
	//
	// [RFC2389]: https://datatracker.ietf.org/doc/html/rfc2389
	// [RFC959]: https://datatracker.ietf.org/doc/html/rfc959
	features map[string]string
}

// NewConfig returns default FTP server configuration options.
func NewConfig(opts ...Option) Config {
	cfg := Config{
		host:        "127.0.0.1",
		readTO:      150 * time.Millisecond,
		writeTO:     50 * time.Millisecond,
		clock:       func() time.Time { return time.Now().UTC() },
		svrReadyMsg: ServerReady,
		features:    make(map[string]string, 20),
	}

	// Minimum features from RFC 959.
	cfg = cfg.EnableFeature(ftpcmd.Minimal...)
	for _, opt := range opts {
		cfg = opt(cfg)
	}
	return cfg
}

// Features returns a sorted slice of enabled features.
func (cfg Config) Features() []string {
	var feats []string
	for k := range cfg.features {
		feats = append(feats, k)
	}
	sort.Strings(feats)
	return feats
}

// HasFeature returns true when the feature is enabled.
func (cfg Config) HasFeature(feat string) bool {
	_, ok := cfg.features[feat]
	return ok
}

// EnableFeature enables a given FTP feature (command). The original instance
// of [Config] is not changed in any way.
func (cfg Config) EnableFeature(feats ...string) Config {
	cfg.features = maps.Clone(cfg.features)
	for _, feat := range feats {
		name, opts := SplitCmdLine(feat)
		if _, ok := cfg.features[feat]; !ok {
			cfg.features[name] = strings.Join(opts, "")
		}
	}
	return cfg
}

// DisableFeature disables the given FTP feature (command). The original
// instance of Config is not changed in any way.
func (cfg Config) DisableFeature(feat string) Config {
	cfg.features = maps.Clone(cfg.features)
	delete(cfg.features, feat)
	return cfg
}

// DisableAllFeatures disables all currently enabled features. The original
// instance of Config is not changed in any way.
func (cfg Config) DisableAllFeatures() Config {
	cfg.features = map[string]string{}
	return cfg
}
