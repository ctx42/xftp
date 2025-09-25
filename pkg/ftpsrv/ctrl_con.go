package ftpsrv

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/textproto"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

// ErrQuit is returned by "QUIT" FTP command handler.
var ErrQuit = errors.New("connection closing after QUIT command")

// CtrlCon represents FTP control connection.
type CtrlCon struct {
	ses       *Session
	tslCfg    *tls.Config
	conn      net.Conn
	proto     *textproto.Conn // Connection wrapped in text protocol.
	log       zerolog.Logger
	quitCh    chan struct{} // When closed forces connection to be closed.
	quitFn    func()        // Called to close control channel.
	listening bool          // The listener goroutine is listening.
	mx        sync.RWMutex  // Guards struct fields.
}

// NewCtrlCon returns a new instance of [CtrlCon].
func NewCtrlCon(ses *Session, conn net.Conn, log zerolog.Logger) *CtrlCon {
	quitCh := make(chan struct{}, 1)
	return &CtrlCon{
		ses:    ses,
		conn:   conn,
		proto:  textproto.NewConn(conn),
		log:    log,
		quitCh: quitCh,
		quitFn: sync.OnceFunc(func() { close(quitCh) }),
	}
}

// WithTLS sets TLS configuration for the session.
func (cc *CtrlCon) WithTLS(tlsCfg *tls.Config) *CtrlCon {
	cc.tslCfg = tlsCfg
	return cc
}

// Listen starts listening for control commands from the client.
func (cc *CtrlCon) Listen() *CtrlCon {
	cc.mx.Lock()
	defer cc.mx.Unlock()
	if cc.listening {
		return cc
	}
	started := make(chan struct{})
	go cc.listen(started)
	<-started // Return only when the goroutine has started.
	cc.listening = true
	return cc
}

// listen starts listening for commands. It should be listening in goroutine.
// Closes the listening channel once it starts.
func (cc *CtrlCon) listen(started chan struct{}) {
	defer func() {
		LogError(nil, cc.log, cc.close(), nil)
		cc.log.Debug().Msg("cc.listen: exiting")
	}()
	cc.log.Debug().Msg("cc.listen: started")
	close(started)

	if err := cc.writeLine(cc.ses.Cfg.svrReadyMsg); err != nil {
		cc.log.Error().Err(err).Send()
	}

	for {
		select {
		case <-cc.quitCh:
			cc.log.Debug().Msg("cc.listen: quitting")
			return
		default:
		}

		_ = cc.conn.SetReadDeadline(time.Now().Add(cc.ses.Cfg.readTO))
		line, err := cc.proto.ReadLine()
		if err != nil {
			var e *net.OpError
			if errors.As(err, &e) && e.Timeout() {
				cc.log.Debug().Msgf("cc.listen: read timeout")
				continue
			}

			switch {
			case strings.Contains(err.Error(), "connection reset by peer"):
				cc.log.Debug().Msgf("cc.listen: connection reset by peer")

			case errors.Is(err, io.EOF):
				cc.log.Debug().Msgf("cc.listen: closed by client")
			}
			return
		}
		cc.log.Debug().Msgf("< %s", line)

		cmd, args := SplitCmdLine(line)
		cmd = strings.ToUpper(cmd)
		cc.mx.Lock()
		if err = cc.handleCommand(cmd, args...); err != nil {
			if errors.Is(err, ErrQuit) {
				cc.mx.Unlock()
				return
			}
			cc.log.Error().Err(err).Send()
			cc.log.Debug().Msgf("cc.listen: handler error")
		}
		cc.mx.Unlock()
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

// Close closes control connection and data connection if it exists and stops
// listening for control commands. No message is sent to the client.
func (cc *CtrlCon) Close() error {
	cc.quitFn()
	return nil
}

// Close closes control connection and data connection if it exists and stops
// listening for control commands. No message is sent to the client.
func (cc *CtrlCon) close() error {
	cc.mx.Lock()
	defer cc.mx.Unlock()
	if !cc.listening {
		return nil
	}
	_ = cc.conn.SetWriteDeadline(time.Now().Add(cc.ses.Cfg.writeTO))
	meta := map[string]any{"action": "cc.close"}
	LogError(nil, cc.log, cc.proto.Close(), meta)
	cc.log.Debug().Msgf("cc.Close")
	return nil
}
