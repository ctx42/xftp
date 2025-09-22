package ftpsrv

import (
	"testing"
	"time"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_WithHost(t *testing.T) {
	// --- Given ---
	cfg := Config{}

	// --- When ---
	have := WithHost("host")(cfg)

	// --- Then ---
	assert.Equal(t, "host", have.host)
}

func Test_WithClock(t *testing.T) {
	// --- Given ---
	want := time.Now
	cfg := Config{}

	// --- When ---
	have := WithClock(want)(cfg)

	// --- Then ---
	assert.Same(t, want, have.clock)
}

func Test_WithTLS(t *testing.T) {
	// --- Given ---
	cert := []byte{1}
	key := []byte{1}
	cfg := Config{}

	// --- When ---
	have := WithTLS(cert, key)(cfg)

	// --- Then ---
	assert.True(t, have.supportTLS)
	assert.Same(t, cert, have.cert)
	assert.Same(t, key, have.key)
}

func Test_WithImpTLS(t *testing.T) {
	// --- Given ---
	cfg := Config{}

	// --- When ---
	have := WithImpTLS(cfg)

	// --- Then ---
	assert.True(t, have.implicitTLS)
}

func Test_WithReadyMsg(t *testing.T) {
	// --- Given ---
	custom := Resp(220, "SESSION:(abc)", "FTP Server ready.")
	cfg := Config{}

	// --- When ---
	have := WithReadyMsg(custom)(cfg)

	// --- Then ---
	assert.Equal(t, Response{}, cfg.svrReadyMsg)
	assert.Equal(t, custom, have.svrReadyMsg)
}

func Test_NewConfig(t *testing.T) {
	t.Run("no options", func(t *testing.T) {
		// --- When ---
		cfg := NewConfig()

		// --- Then ---
		assert.Equal(t, "127.0.0.1", cfg.host)
		assert.Duration(t, "150ms", cfg.readTO)
		assert.Duration(t, "50ms", cfg.writeTO)
		assert.NotNil(t, cfg.clock)
		assert.False(t, cfg.supportTLS)
		assert.False(t, cfg.implicitTLS)
		assert.Equal(t, ServerReady, cfg.svrReadyMsg)
		assert.Nil(t, cfg.cert)
		assert.Nil(t, cfg.key)

		assert.Fields(t, 9, Config{}) // Update the above assertions on fail.
	})

	t.Run("with option", func(t *testing.T) {
		// --- When ---
		cfg := NewConfig(WithHost("1.2.3.4"))

		// --- Then ---
		assert.Equal(t, "1.2.3.4", cfg.host)
	})
}
