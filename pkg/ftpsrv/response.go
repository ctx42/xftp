package ftpsrv

import (
	"fmt"
	"strconv"
	"strings"
)

// Response represents FTP response.
type Response struct {
	code  int      // Response code.
	pad   bool     // Left pad lines in the middle with space.
	lines []string // Response lines.
}

// Resp returns a new instance of [Response].
func Resp(code int, lines ...string) Response {
	return Response{code: code, lines: lines}
}

// With generates string representation of the response and then treaties it as
// a format string with given arguments.
func (r Response) With(args ...any) string {
	return fmt.Sprintf(r.String(), args...)
}

// Code returns repose code.
func (r Response) Code() int { return r.code }

// Pad adds single space padding to lines in between first and last.
func (r Response) Pad() Response {
	r.pad = true
	return r
}

// String returns string representation of the response.
func (r Response) String() string {
	if r.code == 0 {
		return ""
	}
	if len(r.lines) == 0 {
		return fmt.Sprintf("%d ", r.code)
	}
	if len(r.lines) == 1 {
		return fmt.Sprintf("%d %s", r.code, r.lines[0])
	}

	var msg string
	var code, minus, pad, crlf string
	for i, line := range r.lines {
		line = strings.TrimRight(line, "\r\n") // Don't add CRLF twice.
		first := i == 0
		last := i == len(r.lines)-1
		middle := !first && !last

		if first {
			minus = "-"
			code = strconv.Itoa(r.code)
			crlf = "\r\n"
		}
		if middle {
			if line[0] >= '0' && line[0] <= '9' {
				line = " " + line
			} else if r.pad {
				pad = " "
			}
		}
		if last {
			minus = " "
			code = strconv.Itoa(r.code)
			crlf = ""
		}
		msg += fmt.Sprintf("%s%s%s%s%s", code, minus, pad, line, crlf)
		code = ""
		minus = ""
		pad = ""
	}
	return msg
}

// Text returns response in the same format:
//
//	<code> line1
//	line1
//	...
//	line n
func (r Response) Text() string {
	msg := ""
	for _, line := range r.lines {
		msg += strings.TrimRight(line, "\r\n") + "\n"
	}
	return strconv.Itoa(r.code) + " " + strings.TrimRight(msg, "\n")
}
