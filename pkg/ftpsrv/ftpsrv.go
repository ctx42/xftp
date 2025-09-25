package ftpsrv

import (
	"net"

	"github.com/gofrs/uuid/v5"
)

// ISession represents an FTP session.
type ISession interface {
	// ID returns unique UUIDv7 FTP Session ID.
	ID() uuid.UUID
}

// Feature represents an FTP feature (command).
type Feature interface {
	// Name returns the feature name.
	Name() string

	// RequireAuth returns true if the feature requires authentication.
	RequireAuth() bool

	// Execute executes the feature within the given FTP session.
	Execute(ses ISession, args ...string)
}

// ConnInfo represents an FTP connection info.
type ConnInfo interface {
	// LocalAddr returns the local network address, if known.
	LocalAddr() net.Addr

	// RemoteAddr returns the remote network address, if known.
	RemoteAddr() net.Addr
}
