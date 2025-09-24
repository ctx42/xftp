package ftpsrv

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/textproto"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"
)

// ErrQuit is returned by "QUIT" FTP command handler.
var ErrQuit = errors.New("connection closing after QUIT command")

// Session represents single FTP session.
type Session struct {
	id     uuid.UUID
	cfg    Config
	tslCfg *tls.Config
	conn   net.Conn
	proto  *textproto.Conn // Connection wrapped in text protocol.
}

// NewSession returns a new instance of [Session].
func NewSession(cfg Config, conn net.Conn) *Session {
	return &Session{
		id:    uuid.Must(uuid.NewV7()),
		cfg:   cfg,
		conn:  conn,
		proto: textproto.NewConn(conn),
	}
}

// WithTLS sets TLS configuration for the session.
func (ses *Session) WithTLS(tlsCfg *tls.Config) *Session {
	ses.tslCfg = tlsCfg
	return ses
}

// Listen starts listening for control commands from the client.
func (ses *Session) Listen() *Session {
	started := make(chan struct{})
	ses.listen(started)
	<-started
	return ses
}

func (ses *Session) listen(started chan struct{}) {
	close(started)

	log := ses.cfg.log
	if err := ses.writeLine(ses.cfg.svrReadyMsg); err != nil {
		log.Error().Err(err).Send()
	}

	for {
		_ = ses.conn.SetReadDeadline(time.Now().Add(ses.cfg.readTO))
		line, err := ses.proto.ReadLine()
		if err != nil {
			var e *net.OpError
			if errors.As(err, &e) && e.Timeout() {
				log.Debug().Msgf("SrvCC.listen: read timeout")
				continue
			}

			switch {
			case strings.Contains(err.Error(), "connection reset by peer"):
				log.Debug().Msgf("SrvCC.listen: connection reset by peer")

			case errors.Is(err, io.EOF):
				log.Debug().Msgf("SrvCC.listen: closed by the client")
			}
			return
		}
		log.Debug().Msgf("< %s", line)

		cmd, args := SplitCmdLine(line)
		cmd = strings.ToUpper(cmd)
		if err = ses.handleCommand(cmd, args...); err != nil {
			if errors.Is(err, ErrQuit) {
				return
			}
			log.Error().Err(err).Send()
			log.Debug().Msgf("SrvCC.listen: handler error")
		}
	}
}

// writeLine formats and writes the given arguments to the [textproto.Writer]
// using the specified format string.
func (ses *Session) writeLine(resp Response, args ...any) error {
	log := ses.cfg.log
	msg := resp.With(args...)
	if msg == "" {
		log.Debug().Msg(">")
		return nil
	}

	log.Debug().Msgf("> %s", msg)
	_ = ses.conn.SetWriteDeadline(time.Now().Add(ses.cfg.writeTO))
	if err := ses.proto.PrintfLine("%s", msg); err != nil {
		return fmt.Errorf("SrvCC.writeLine: line send: %w, msg: %s", err, msg)
	}
	return nil
}

func (ses *Session) handleCommand(cmd string, args ...string) error {
	return ses.writeLine(ErrorUnkCmd, cmd)
}
