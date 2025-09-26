package ftpsrv

import (
	"testing"

	"github.com/ctx42/logkit/pkg/logkit"
	"github.com/ctx42/testing/pkg/tester"
	"github.com/gofrs/uuid/v5"
	"github.com/rs/zerolog"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) { goleak.VerifyTestMain(m) }

func TstLogger(t tester.T) (*logkit.Tester, zerolog.Logger) {
	// TODO(rz): test this.
	// TODO(rz): document this.
	tlog := logkit.New(t)
	zlog := zerolog.New(tlog)
	return tlog, zlog
}

func TstID() uuid.UUID {
	// TODO(rz): test this.
	// TODO(rz): document this.
	return uuid.Must(uuid.NewV7())
}
