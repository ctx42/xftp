package ftpsrv

import (
	"net"
	"time"

	"github.com/gofrs/uuid/v5"
)

// Session represents single FTP session.
type Session struct {
	id         uuid.UUID
	cfg        Config
	localAddr  net.Addr
	remoteAddr net.Addr
	startedAt  time.Time
	quitAt     time.Time
	closedAt   time.Time
}

// NewSession returns a new instance of [Session].
func NewSession(cfg Config, local, remote net.Addr) *Session {
	return &Session{
		id:         uuid.Must(uuid.NewV7()),
		cfg:        cfg,
		startedAt:  cfg.clock(),
		localAddr:  local,
		remoteAddr: remote,
	}
}
