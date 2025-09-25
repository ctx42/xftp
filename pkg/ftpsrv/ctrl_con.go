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

	"github.com/rs/zerolog"
)

// ErrQuit is returned by "QUIT" FTP command handler.
var ErrQuit = errors.New("connection closing after QUIT command")

// CtrlCon represents FTP control connection.
type CtrlCon struct {
	ses    *Session
	tslCfg *tls.Config
	conn   net.Conn
	proto  *textproto.Conn // Connection wrapped in text protocol.
	log    zerolog.Logger
}

// NewCtrlCon returns a new instance of [CtrlCon].
func NewCtrlCon(ses *Session, conn net.Conn, log zerolog.Logger) *CtrlCon {
	return &CtrlCon{
		ses:   ses,
		conn:  conn,
		proto: textproto.NewConn(conn),
		log:   log,
	}
}

// WithTLS sets TLS configuration for the session.
func (cc *CtrlCon) WithTLS(tlsCfg *tls.Config) *CtrlCon {
	cc.tslCfg = tlsCfg
	return cc
}

// Listen starts listening for control commands from the client.
func (cc *CtrlCon) Listen() *CtrlCon {
	started := make(chan struct{})
	go cc.listen(started)
	<-started
	return cc
}

func (cc *CtrlCon) listen(started chan struct{}) {
	close(started)

	log := cc.log
	if err := cc.writeLine(cc.ses.Cfg.svrReadyMsg); err != nil {
		log.Error().Err(err).Send()
	}

	for {
		_ = cc.conn.SetReadDeadline(time.Now().Add(cc.ses.Cfg.readTO))
		line, err := cc.proto.ReadLine()
		if err != nil {
			var e *net.OpError
			if errors.As(err, &e) && e.Timeout() {
				log.Debug().Msg("cc.listen: read timeout")
				continue
			}

			switch {
			case strings.Contains(err.Error(), "connection reset by peer"):
				log.Debug().Msg("cc.listen: connection reset by peer")

			case errors.Is(err, io.EOF):
				log.Debug().Msg("cc.listen: closed by the client")
			}
			return
		}
		log.Debug().Msgf("< %s", line)

		cmd, args := SplitCmdLine(line)
		cmd = strings.ToUpper(cmd)
		if err = cc.handleCommand(cmd, args...); err != nil {
			if errors.Is(err, ErrQuit) {
				return
			}
			log.Error().Err(err).Send()
			log.Debug().Msg("cc.listen: handler error")
		}
	}
}

// writeLine formats and writes the given arguments to the [textproto.Writer]
// using the specified format string.
func (cc *CtrlCon) writeLine(resp Response, args ...any) error {
	msg := resp.With(args...)
	if msg == "" {
		cc.log.Debug().Msg(">")
		return nil
	}

	cc.log.Debug().Msgf("> %s", msg)
	_ = cc.conn.SetWriteDeadline(time.Now().Add(cc.ses.Cfg.writeTO))
	if err := cc.proto.PrintfLine("%s", msg); err != nil {
		return fmt.Errorf("cc.writeLine: line send: %w, msg: %s", err, msg)
	}
	return nil
}

func (cc *CtrlCon) handleCommand(cmd string, args ...string) error {
	return cc.writeLine(ErrorUnkCmd, cmd)
}
