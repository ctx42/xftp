package ftpsrv

import (
	"testing"
	"time"

	"github.com/ctx42/testing/pkg/assert"

	"github.com/ctx42/xftp/pkg/ftpcmd"
)

// TODO(rz): Test WithoutFeature, WithFeature, WithHost

func Test_WithHost(t *testing.T) {
	// --- Given ---
	cfg := Config{}

	// --- When ---
	have := WithHost("host")(cfg)

	// --- Then ---
	assert.Empty(t, cfg.host)
	assert.Equal(t, "host", have.host)
}

func Test_WithClock(t *testing.T) {
	// --- Given ---
	cfg := Config{}
	want := time.Now

	// --- When ---
	have := WithClock(want)(cfg)

	// --- Then ---
	assert.NotSame(t, cfg.clock, have.clock)
	assert.Same(t, want, have.clock)
}

func Test_WithTLS(t *testing.T) {
	// --- Given ---
	cfg := Config{features: make(map[string]string)}

	// --- When ---
	have := WithTLS(cfg)

	// --- Then ---
	assert.False(t, cfg.supportTLS)
	assert.True(t, have.supportTLS)
	wFeatures := map[string]string{"AUTH": "TLS", "PBSZ": "", "PROT": ""}
	assert.Equal(t, wFeatures, have.features)
}

func Test_WithImpTLS(t *testing.T) {
	// --- Given ---
	cfg := Config{features: make(map[string]string)}

	// --- When ---
	have := WithImpTLS(cfg)

	// --- Then ---
	assert.False(t, cfg.supportTLS)
	assert.True(t, have.supportTLS)
	assert.Equal(t, "TLS", have.features[ftpcmd.AUTH])
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
		want := map[string]string{
			ftpcmd.MODE: "",
			ftpcmd.NOOP: "",
			ftpcmd.QUIT: "",
			ftpcmd.RETR: "",
			ftpcmd.STOR: "",
			ftpcmd.STRU: "",
			ftpcmd.TYPE: "",
			ftpcmd.USER: "",
		}
		assert.Equal(t, want, cfg.features)
		assert.Equal(t, ServerReady, cfg.svrReadyMsg)

		assert.Fields(t, 8, Config{}) // Update the above assertions on fail.
	})

	t.Run("with option", func(t *testing.T) {
		// --- When ---
		cfg := NewConfig(WithFeature(ftpcmd.MDTM))

		// --- Then ---
		assert.Equal(t, "127.0.0.1", cfg.host)
		assert.Duration(t, "150ms", cfg.readTO)
		assert.Duration(t, "50ms", cfg.writeTO)
		assert.NotNil(t, cfg.clock)
		assert.False(t, cfg.supportTLS)
		want := map[string]string{
			ftpcmd.MDTM: "",
			ftpcmd.MODE: "",
			ftpcmd.NOOP: "",
			ftpcmd.QUIT: "",
			ftpcmd.RETR: "",
			ftpcmd.STOR: "",
			ftpcmd.STRU: "",
			ftpcmd.TYPE: "",
			ftpcmd.USER: "",
		}
		assert.Equal(t, want, cfg.features)
		assert.Equal(t, ServerReady, cfg.svrReadyMsg)

		assert.Fields(t, 8, Config{}) // Update the above assertions on fail.
	})
}

func Test_Config_Features(t *testing.T) {
	// --- Given ---
	cfg := NewConfig()

	// --- When ---
	have := cfg.Features()

	// --- Then ---
	want := []string{
		ftpcmd.MODE,
		ftpcmd.NOOP,
		ftpcmd.QUIT,
		ftpcmd.RETR,
		ftpcmd.STOR,
		ftpcmd.STRU,
		ftpcmd.TYPE,
		ftpcmd.USER,
	}
	assert.Equal(t, want, have)
}

func Test_Config_HasFeature(t *testing.T) {
	t.Run("has", func(t *testing.T) {
		// --- Given ---
		cfg := NewConfig()

		// --- When ---
		have := cfg.HasFeature(ftpcmd.USER)

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("does not have", func(t *testing.T) {
		// --- Given ---
		cfg := NewConfig()

		// --- When ---
		have := cfg.HasFeature("ABCD")

		// --- Then ---
		assert.False(t, have)
	})
}

func Test_Config_EnableFeature(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		cfg := NewConfig()

		// --- When ---
		have := cfg.EnableFeature("APX")

		// --- Then ---
		want := map[string]string{
			ftpcmd.MODE: "",
			ftpcmd.NOOP: "",
			ftpcmd.QUIT: "",
			ftpcmd.RETR: "",
			ftpcmd.STOR: "",
			ftpcmd.STRU: "",
			ftpcmd.TYPE: "",
			ftpcmd.USER: "",
		}
		assert.Equal(t, want, cfg.features)

		want = map[string]string{
			"APX":       "",
			ftpcmd.MODE: "",
			ftpcmd.NOOP: "",
			ftpcmd.QUIT: "",
			ftpcmd.RETR: "",
			ftpcmd.STOR: "",
			ftpcmd.STRU: "",
			ftpcmd.TYPE: "",
			ftpcmd.USER: "",
		}
		assert.Equal(t, want, have.features)
	})

	t.Run("enable existing", func(t *testing.T) {
		// --- Given ---
		cfg := NewConfig()

		// --- When ---
		have := cfg.EnableFeature(ftpcmd.MODE)

		// --- Then ---
		want := map[string]string{
			ftpcmd.MODE: "",
			ftpcmd.NOOP: "",
			ftpcmd.QUIT: "",
			ftpcmd.RETR: "",
			ftpcmd.STOR: "",
			ftpcmd.STRU: "",
			ftpcmd.TYPE: "",
			ftpcmd.USER: "",
		}
		assert.Equal(t, want, cfg.features)
		assert.Equal(t, want, have.features)
	})

	t.Run("enable with options", func(t *testing.T) {
		// --- Given ---
		cfg := NewConfig()

		// --- When ---
		have := cfg.EnableFeature(ftpcmd.MLST + " type*;size*;modify*;")

		// --- Then ---
		want := map[string]string{
			ftpcmd.MODE: "",
			ftpcmd.NOOP: "",
			ftpcmd.QUIT: "",
			ftpcmd.RETR: "",
			ftpcmd.STOR: "",
			ftpcmd.STRU: "",
			ftpcmd.TYPE: "",
			ftpcmd.USER: "",
		}
		assert.Equal(t, want, cfg.features)

		want = map[string]string{
			ftpcmd.MLST: "type*;size*;modify*;",
			ftpcmd.MODE: "",
			ftpcmd.NOOP: "",
			ftpcmd.QUIT: "",
			ftpcmd.RETR: "",
			ftpcmd.STOR: "",
			ftpcmd.STRU: "",
			ftpcmd.TYPE: "",
			ftpcmd.USER: "",
		}
		assert.Equal(t, want, have.features)
	})
}

func Test_Config_DisableFeature(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		cfg := NewConfig()

		// --- When ---
		have := cfg.DisableFeature(ftpcmd.MODE)

		// --- Then ---
		want := map[string]string{
			ftpcmd.MODE: "",
			ftpcmd.NOOP: "",
			ftpcmd.QUIT: "",
			ftpcmd.RETR: "",
			ftpcmd.STOR: "",
			ftpcmd.STRU: "",
			ftpcmd.TYPE: "",
			ftpcmd.USER: "",
		}
		assert.Equal(t, want, cfg.features)

		want = map[string]string{
			ftpcmd.NOOP: "",
			ftpcmd.QUIT: "",
			ftpcmd.RETR: "",
			ftpcmd.STOR: "",
			ftpcmd.STRU: "",
			ftpcmd.TYPE: "",
			ftpcmd.USER: "",
		}
		assert.Equal(t, want, have.features)
	})

	t.Run("delete not exising", func(t *testing.T) {
		// --- Given ---
		cfg := NewConfig()

		// --- When ---
		have := cfg.DisableFeature("APX")

		// --- Then ---
		want := map[string]string{
			ftpcmd.MODE: "",
			ftpcmd.NOOP: "",
			ftpcmd.QUIT: "",
			ftpcmd.RETR: "",
			ftpcmd.STOR: "",
			ftpcmd.STRU: "",
			ftpcmd.TYPE: "",
			ftpcmd.USER: "",
		}
		assert.Equal(t, want, cfg.features)
		assert.Equal(t, want, have.features)
	})
}

func Test_Config_WithoutFeature(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		cfg := NewConfig()

		// --- When ---
		have := WithoutFeature(ftpcmd.MODE)(cfg)

		// --- Then ---
		want := map[string]string{
			ftpcmd.MODE: "",
			ftpcmd.NOOP: "",
			ftpcmd.QUIT: "",
			ftpcmd.RETR: "",
			ftpcmd.STOR: "",
			ftpcmd.STRU: "",
			ftpcmd.TYPE: "",
			ftpcmd.USER: "",
		}
		assert.Equal(t, want, cfg.features)

		want = map[string]string{
			ftpcmd.NOOP: "",
			ftpcmd.QUIT: "",
			ftpcmd.RETR: "",
			ftpcmd.STOR: "",
			ftpcmd.STRU: "",
			ftpcmd.TYPE: "",
			ftpcmd.USER: "",
		}
		assert.Equal(t, want, have.features)
	})

	t.Run("delete not exising", func(t *testing.T) {
		// --- Given ---
		cfg := NewConfig()

		// --- When ---
		have := WithoutFeature("APX")(cfg)

		// --- Then ---
		want := map[string]string{
			ftpcmd.MODE: "",
			ftpcmd.NOOP: "",
			ftpcmd.QUIT: "",
			ftpcmd.RETR: "",
			ftpcmd.STOR: "",
			ftpcmd.STRU: "",
			ftpcmd.TYPE: "",
			ftpcmd.USER: "",
		}
		assert.Equal(t, want, cfg.features)
		assert.Equal(t, want, have.features)
	})
}

func Test_Config_DisableAllFeatures(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		cfg := NewConfig()

		// --- When ---
		have := cfg.DisableAllFeatures()

		// --- Then ---
		want := map[string]string{
			ftpcmd.MODE: "",
			ftpcmd.NOOP: "",
			ftpcmd.QUIT: "",
			ftpcmd.RETR: "",
			ftpcmd.STOR: "",
			ftpcmd.STRU: "",
			ftpcmd.TYPE: "",
			ftpcmd.USER: "",
		}
		assert.Equal(t, want, cfg.features)
		assert.Empty(t, have.features)
	})
}
