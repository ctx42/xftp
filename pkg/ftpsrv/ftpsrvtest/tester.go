package ftpsrvtest

import (
	"fmt"
	"io"
	"net"
	"net/textproto"
	"strconv"
	"time"

	"github.com/ctx42/logkit/pkg/logkit"
	"github.com/ctx42/testing/pkg/kit/timekit"
	"github.com/ctx42/testing/pkg/notice"
	"github.com/ctx42/testing/pkg/tester"
	"github.com/rs/zerolog"

	"github.com/ctx42/xftp/pkg/ftpsrv"
)

type hidLogTrait = logkit.Trait // Don't export embedded field.

// Tester is a helper for testing FTP server.
type Tester struct {
	*hidLogTrait                  // Log test helper.
	cliCC        net.Conn         // Client side of the control connection.
	srvCC        net.Conn         // Server side of the control connection.
	rto          time.Duration    // Timeout reading from control connection.
	wto          time.Duration    // Timeout writing to control connection.
	clk          func() time.Time // Clock to use (2000-01-01 00:00:00 UTC).
	log          zerolog.Logger   // Test logger.
	t            tester.T         // Test manager.
}

// NewTester returns a new instance of [Tester].
func NewTester(t tester.T) *Tester {
	t.Helper()
	tlog := logkit.NewTrait(t)
	now := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	return &Tester{
		hidLogTrait: tlog,
		log:         zerolog.New(tlog.LogWriter()),
		rto:         50 * time.Millisecond,
		wto:         50 * time.Millisecond,
		clk:         timekit.ClockStartingAt(now),
		t:           t,
	}
}

// INE doesn't mark the test as failed when the logs weren't examined, and
// there are no log messages with error or panic log levels.
func (tst *Tester) INE() *Tester {
	tst.t.Helper()
	tst.IgnoreNonErrorLogs()
	return tst
}

// WireUp creates server and client side [net.Conn] instance and registers
// [Tester.cleanup] to be called after the test ends. Terminates the test on
// error.
func (tst *Tester) WireUp() *Tester {
	tst.t.Helper()
	if tst.srvCC != nil {
		msg := "tester: server control connection can be created only once"
		tst.t.Fatal(msg)
		return tst
	}
	if tst.cliCC != nil {
		msg := "tester: client control connection can be created only once"
		tst.t.Fatal(msg)
		return tst
	}
	tst.srvCC, tst.cliCC = net.Pipe()
	tst.t.Cleanup(tst.cleanup)
	return tst
}

// Logger returns test logger.
func (tst *Tester) Logger() zerolog.Logger {
	tst.t.Helper()
	return tst.log
}

// Clock returns clock used in the [Tester].
func (tst *Tester) Clock() func() time.Time { return tst.clk }

// SrvCon returns server side control connection created by [Tester.WireUp].
func (tst *Tester) SrvCon() net.Conn { return tst.srvCC }

// CliConClose closes client side of control connection without sending any
// commands. If closing fails, it marks the test as failed with an appropriate
// error message.
func (tst *Tester) CliConClose() *Tester {
	tst.t.Helper()
	if err := tst.cliCC.Close(); err != nil {
		tst.t.Errorf("tester.cli_con_close: %s", err)
		tst.cliCC = nil
	}
	return tst
}

// Config returns [ftpsrv.Config] instance with values adjusted for testing.
func (tst *Tester) Config(opts ...ftpsrv.Option) ftpsrv.Config {
	tst.t.Helper()
	defaults := []ftpsrv.Option{
		ftpsrv.WithReadTimeout(tst.rto),
		ftpsrv.WithWriteTimeout(tst.wto),
		ftpsrv.WithClock(tst.clk),
	}
	return ftpsrv.NewConfig(append(defaults, opts...)...)
}

// Session returns [ftpsrv.Session] instance with values reflecting the tester.
func (tst *Tester) Session(opts ...ftpsrv.Option) *ftpsrv.Session {
	tst.t.Helper()
	return ftpsrv.NewSession(
		tst.Config(opts...),
		tst.srvCC.LocalAddr(),
		tst.srvCC.RemoteAddr(),
	)
}

// SendCmd sends the command through the control connection. If the sending
// fails, it marks the test as failed with an appropriate error message.
func (tst *Tester) SendCmd(cmd string, format string, args ...any) *Tester {
	tst.t.Helper()
	if opt := fmt.Sprintf(format, args...); opt != "" {
		cmd += " " + opt
	}
	if err := tst.writeLine("%s", cmd); err != nil {
		tst.t.Errorf("tester.send_cmd: %s; tester.write_line: %s", cmd, err)
	}
	return tst
}

// GetReply reads the server side control connection and asserts it's equal to
// the msg. When getting the replay fails, or the message does not match, it
// marks the test as failed with an appropriate error message.
func (tst *Tester) GetReply(msg string) *Tester {
	tst.t.Helper()
	tst.GetReplyf("%s", msg)
	return tst
}

// GetReplyf reads the server side control connection and asserts it's equal to
// the formated message. When getting the replay fails, or the message does not
// match, it marks the test as failed with an appropriate error message.
func (tst *Tester) GetReplyf(format string, args ...any) *Tester {
	tst.t.Helper()
	have, err := tst.readLine()
	if err != nil {
		tst.t.Errorf("tester.get_reply: Tester.readLine: %s", err)
		return tst
	}
	want := fmt.Sprintf(format, args...)
	if want != have {
		msg := notice.New("expected server replay to equal").
			Want("%s", want).
			Have("%s", have)
		tst.t.Error(msg)
	}
	return tst
}

// readLine reads lines from the client side of the control connection. When
// reading fails, it marks the test as failed with an appropriate error message.
func (tst *Tester) readLine() (string, error) {
	tst.t.Helper()
	_ = tst.cliCC.SetReadDeadline(time.Now().Add(tst.rto))
	conn := textproto.NewConn(tst.cliCC)
	code, msg, err := conn.ReadResponse(-1)
	if err != nil {
		return "", err
	}
	line := strconv.Itoa(code) + " " + msg
	return line, nil
}

// writeLine writes line to the client side of the control connection.
func (tst *Tester) writeLine(format string, args ...any) error {
	tst.t.Helper()
	_ = tst.cliCC.SetWriteDeadline(time.Now().Add(tst.wto))
	w := textproto.NewConn(tst.cliCC).Writer
	return w.PrintfLine(format, args...)
}

// CloseAfterTest registers the c.Close method to be called after the test ends.
func (tst *Tester) CloseAfterTest(c io.Closer) {
	tst.t.Helper()
	tst.t.Cleanup(func() {
		tst.t.Helper()
		if err := c.Close(); err != nil {
			tst.t.Errorf("tester.close_after_test: %s", err)
		}
	})
}

// cleanup closes all connections and releases the tester resources.
func (tst *Tester) cleanup() {
	tst.t.Helper()
	if tst.srvCC != nil {
		if err := tst.srvCC.Close(); err != nil {
			tst.t.Errorf("tester.cleanup.srv: %s", err)
		}
	}
	if tst.cliCC != nil {
		if err := tst.cliCC.Close(); err != nil {
			tst.t.Errorf("tester.cleanup.cli: %s", err)
		}
	}
}
