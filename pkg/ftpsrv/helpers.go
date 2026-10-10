package ftpsrv

import (
	"slices"
	"strings"
)

// SplitCmdLine splits FTP command to command line and its arguments.
func SplitCmdLine(line string) (string, []string) {
	var cmd string
	args := strings.Split(line, " ")
	switch len(args) {
	case 0:
		cmd = ""
	case 1:
		cmd = args[0]
		args = args[:0]
	default:
		cmd = args[0]
		args = args[1:]
	}
	args = slices.DeleteFunc(args, func(s string) bool {
		return strings.TrimSpace(s) == ""
	})
	if len(args) == 0 {
		args = nil
	}
	return cmd, args
}
