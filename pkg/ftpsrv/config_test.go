package ftpsrv

import (
	"testing"
	"time"

	"github.com/ctx42/testing/pkg/assert"

	"github.com/ctx42/xftp/pkg/ftpcmd"
)

func Test_WithHost(t *testing.T) {
	// --- Given ---
	cfg := Config{}

	// --- When ---
	have := WithHost("host")(cfg)

	// --- Then ---
	assert.Equal(t, "host", have.host)
}

func Test_WithPort(t *testing.T) {
	// --- Given ---
	cfg := Config{}

	// --- When ---
	have := WithPort(2121)(cfg)

	// --- Then ---
	assert.Equal(t, 2121, have.port)
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

func Test_WithReadTimeout(t *testing.T) {
	// --- Given ---
	cfg := Config{}

	// --- When ---
	have := WithReadTimeout(time.Second)(cfg)

	// --- Then ---
	assert.Equal(t, time.Second, have.readTO)
}

func Test_WithWriteTimeout(t *testing.T) {
	// --- Given ---
	cfg := Config{}

	// --- When ---
	have := WithWriteTimeout(time.Second)(cfg)

	// --- Then ---
	assert.Equal(t, time.Second, have.writeTO)
}

func Test_WithoutFeature(t *testing.T) {
	// --- Given ---
	cfg := Config{features: map[string]Command{"XXXX": nil}}

	// --- When ---
	have := WithoutFeature("XXXX")(cfg)

	// --- Then ---
	assert.Len(t, 0, have.features)
}

func Test_NewConfig(t *testing.T) {
	t.Run("no options", func(t *testing.T) {
		// --- When ---
		cfg := NewConfig()

		// --- Then ---
		assert.Equal(t, "127.0.0.1", cfg.host)
		assert.Equal(t, 21, cfg.port)
		assert.Duration(t, "150ms", cfg.readTO)
		assert.Duration(t, "50ms", cfg.writeTO)
		assert.NotNil(t, cfg.clock)
		assert.False(t, cfg.supportTLS)
		assert.False(t, cfg.implicitTLS)
		assert.Equal(t, ServerReady, cfg.svrReadyMsg)
		assert.Nil(t, cfg.cert)
		assert.Nil(t, cfg.key)
		assert.Len(t, 1, cfg.features)

		assert.Fields(t, 11, Config{}) // Update the above assertions on fail.
	})

	t.Run("with option", func(t *testing.T) {
		// --- When ---
		cfg := NewConfig(WithHost("1.2.3.4"))

		// --- Then ---
		assert.Equal(t, "1.2.3.4", cfg.host)
	})
}

func Test_Config_Features(t *testing.T) {
	// --- Given ---
	cfg := Config{features: map[string]Command{ftpcmd.NOOP: nil}}

	// --- When ---
	have := cfg.Features()

	// --- Then ---
	want := []string{ftpcmd.NOOP}
	assert.Equal(t, want, have)
}

func Test_Config_HasFeature(t *testing.T) {
	t.Run("existing", func(t *testing.T) {
		// --- Given ---
		cfg := Config{features: map[string]Command{ftpcmd.NOOP: nil}}

		// --- When ---
		have := cfg.HasFeature(ftpcmd.NOOP)

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("not existing", func(t *testing.T) {
		// --- Given ---
		cfg := Config{features: map[string]Command{ftpcmd.NOOP: nil}}

		// --- When ---
		have := cfg.HasFeature("XXXX")

		// --- Then ---
		assert.False(t, have)
	})
}

func Test_Config_DisableFeature(t *testing.T) {
	t.Run("existing", func(t *testing.T) {
		// --- Given ---
		cfg := Config{features: map[string]Command{ftpcmd.NOOP: nil}}

		// --- When ---
		have := cfg.DisableFeature(ftpcmd.NOOP)

		// --- Then ---
		assert.HasKey(t, ftpcmd.NOOP, cfg.features)
		assert.HasNoKey(t, ftpcmd.NOOP, have.features)
	})

	t.Run("not existing", func(t *testing.T) {
		// --- Given ---
		cfg := Config{features: map[string]Command{ftpcmd.NOOP: nil}}

		// --- When ---
		have := cfg.DisableFeature("XXXX")

		// --- Then ---
		assert.HasKey(t, ftpcmd.NOOP, cfg.features)
		assert.HasKey(t, ftpcmd.NOOP, have.features)
	})
}

func Test_Config_DisableAllFeatures(t *testing.T) {
	t.Run("existing", func(t *testing.T) {
		// --- Given ---
		cfg := Config{
			features: map[string]Command{
				ftpcmd.NOOP: nil,
				"XXXX":      nil,
			},
		}

		// --- When ---
		have := cfg.DisableAllFeatures()

		// --- Then ---
		assert.Len(t, 2, cfg.features)
		assert.Len(t, 0, have.features)
	})
}
