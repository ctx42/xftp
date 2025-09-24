package ftpsrv

import (
	"github.com/gofrs/uuid/v5"
)

// Session represents single FTP session.
type Session struct {
	id  uuid.UUID
	cfg Config
}

// NewSession returns a new instance of [Session].
func NewSession(id uuid.UUID, cfg Config) *Session {
	return &Session{
		id:  id,
		cfg: cfg,
	}
}
