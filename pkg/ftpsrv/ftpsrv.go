package ftpsrv

import (
	"github.com/gofrs/uuid/v5"
)

// Flags represent an FTP command metadata bits.
type Flags uint32

// Command flags.
const (
	FlagBase     Flags = 1 // Command is defined in RFC959.
	FlagMinimal  Flags = 2 // Command is defined in RFC959 as a minimal requirement.
	FlagExtended Flags = 4 // Command is defined in an RFC.
)

// Command represents an FTP command.
type Command interface {
	// Name returns the command name as it appears in the control connection.
	Name() string

	// Flags returns command kind.
	Flags() Flags

	// Handle handles the command.
	Handle(cc *CtrlCon, args ...string) error
}

// ISession represents an FTP session.
type ISession interface {
	// ID returns unique UUIDv7 FTP Session ID.
	ID() uuid.UUID
}

// type Session interface {
//
// }
