package ftpsrv

import (
	"github.com/ctx42/xftp/pkg/ftpcmd"
)

// handlers is a map of all supported FTP command handlers.
var handlers = map[string]Command{
	ftpcmd.NOOP: CmdNOOP{},
}
