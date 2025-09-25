package ftpsrv

import (
	"net"
	"time"

	"github.com/gofrs/uuid/v5"
)

// Session represents single FTP session.
type Session struct {
	ID         uuid.UUID
	Cfg        Config
	LocalAddr  net.Addr
	RemoteAddr net.Addr
	StartedAt  time.Time
	QuitAt     time.Time
	ClosedAt   time.Time
}
