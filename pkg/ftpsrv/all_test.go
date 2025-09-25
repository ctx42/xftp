package ftpsrv

import (
	"github.com/ctx42/logkit/pkg/logkit"
	"github.com/ctx42/testing/pkg/tester"
	"github.com/gofrs/uuid/v5"
	"github.com/rs/zerolog"
)

func TstLogger(t tester.T) (*logkit.Tester, zerolog.Logger) {
	tlog := logkit.New(t)
	zlog := zerolog.New(tlog)
	return tlog, zlog
}

func TstID() uuid.UUID { return uuid.Must(uuid.NewV7()) }
