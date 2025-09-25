package ftpsrvtest

import (
	"fmt"
	"net"
	"net/textproto"
	"strconv"
	"time"

	"github.com/ctx42/logkit/pkg/logkit"
	"github.com/ctx42/testing/pkg/kit/timekit"
	"github.com/ctx42/testing/pkg/notice"
	"github.com/ctx42/testing/pkg/tester"
	"github.com/gofrs/uuid/v5"
	"github.com/rs/zerolog"

	"github.com/ctx42/xftp/pkg/ftpsrv"
)

type hidLogTrait = logkit.Trait // Don't export embedded field.

type Tester struct {
	*hidLogTrait               // Log test helper.
	cliCC        net.Conn      // Client side of the control connection.
	cliDC        net.Conn      // Client side of a data connection.
	srvCC        net.Conn      // Server side of the control connection.
	replays      []string      // Lines read from the client side control connection.
	rTO          time.Duration // Timeout reading from control connection.
	wTO          time.Duration // Timeout writing to control connection.
	clk          func() time.Time
	log          zerolog.Logger
	t            tester.T
}

func NewTester(t tester.T) *Tester {
	t.Helper()
	now := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
	t.Helper()
	tst := &Tester{
		hidLogTrait: logkit.NewTrait(t),
		replays:     make([]string, 0, 10),
		rTO:         50 * time.Millisecond,
		wTO:         50 * time.Millisecond,
		clk:         timekit.ClockStartingAt(now),
		t:           t,
	}
	tst.log = zerolog.New(tst.LogWriter())
	return tst
}

func (tst *Tester) WireUp() *Tester {
	tst.t.Helper()
	if tst.srvCC != nil {
		msg := "Tester.WireUp: can be created only once"
		tst.t.Fatal(msg)
		return tst
	}
	if tst.cliCC != nil {
		msg := "Tester.WireUp: can be created only once"
		tst.t.Fatal(msg)
		return tst
	}
	tst.srvCC, tst.cliCC = net.Pipe()
	return tst
}

func (tst *Tester) Logger() zerolog.Logger {
	tst.t.Helper()
	return tst.log
}

func (tst *Tester) SrvCon() net.Conn  { return tst.srvCC }
func (tst *Tester) DataCon() net.Conn { return tst.cliDC }

func (tst *Tester) TstConfig(opts ...ftpsrv.Option) ftpsrv.Config {
	tst.t.Helper()
	defaults := []ftpsrv.Option{
		ftpsrv.WithReadTimeout(tst.rTO),
		ftpsrv.WithWriteTimeout(tst.wTO),
		ftpsrv.WithClock(tst.clk),
	}
	return ftpsrv.NewConfig(append(defaults, opts...)...)
}

func (tst *Tester) TstSession(opts ...ftpsrv.Option) *ftpsrv.Session {
	tst.t.Helper()
	ses := &ftpsrv.Session{
		ID:         uuid.Must(uuid.NewV7()),
		Cfg:        tst.TstConfig(opts...),
		LocalAddr:  tst.srvCC.LocalAddr(),
		RemoteAddr: tst.cliCC.RemoteAddr(),
		StartedAt:  tst.clk(),
	}
	return ses
}

// Reply reads server side control connection and asserts it's equal to msg.
// When reply doesn't match it marks the test as failed and writes error
// message to test log.
func (tst *Tester) Reply(msg string) *Tester {
	tst.t.Helper()
	tst.Replyf("%s", msg)
	return tst
}

// Replyf reads server side control connection and asserts it's equal to
// formated message. When reply doesn't match it marks the test as failed and
// writes error message to test log.
func (tst *Tester) Replyf(format string, args ...any) *Tester {
	tst.t.Helper()
	have, err := tst.readLine()
	if err != nil {
		tst.t.Errorf("Tester.Reply: Tester.readLine: %s", err)
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

// readLine reads and returns lines read from client side control connection.
//
// On error, marks the test as failed and writes error message to test log, and
// returns empty string.
func (tst *Tester) readLine() (string, error) {
	tst.t.Helper()
	_ = tst.cliCC.SetReadDeadline(time.Now().Add(tst.rTO))
	conn := textproto.NewConn(tst.cliCC)
	code, msg, err := conn.ReadResponse(-1)
	if err != nil {
		return "", err
	}
	line := strconv.Itoa(code) + " " + msg
	tst.replays = append(tst.replays, line)
	return line, nil
}
