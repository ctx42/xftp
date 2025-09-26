package ftpsrv

import (
	"net"
	"time"

	"github.com/gofrs/uuid/v5"
)

// Session represents single FTP session.
type Session struct {
	id         uuid.UUID
	localAddr  net.Addr
	remoteAddr net.Addr
	startedAt  time.Time
	quitAt     time.Time
	closedAt   time.Time
}

// NewSession returns a new instance of [Session].
func NewSession(startedAt time.Time, local, remote net.Addr) *Session {
	return &Session{
		id:         uuid.Must(uuid.NewV7()),
		startedAt:  startedAt,
		localAddr:  local,
		remoteAddr: remote,
	}
}
func (ses *Session) ID() uuid.UUID { return ses.id }
