package ftpsrv

import (
	"github.com/gofrs/uuid/v5"
)

// Session represents an FTP session.
type Session interface {
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
	Execute(ses Session, args ...string)
}
