package ftpsrvtest

import (
	"net"
	"time"

	"github.com/ctx42/testing/pkg/kit/timekit"
	"github.com/ctx42/testing/pkg/tester"
	"github.com/gofrs/uuid/v5"

	"github.com/ctx42/xftp/pkg/ftpsrv"
)

type Tester struct {
	cliCC   net.Conn      // Client side of the control connection.
	cliDC   net.Conn      // Client side of a data connection.
	srvCC   net.Conn      // Server side of the control connection.
	replays []string      // Lines read from the client side control connection.
	rTO     time.Duration // Timeout reading from control connection.
	wTO     time.Duration // Timeout writing to control connection.
	clk     func() time.Time
	t       tester.T
}

func NewTester(t tester.T) *Tester {
	now := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
	t.Helper()
	tst := &Tester{
		replays: make([]string, 0, 10),
		rTO:     50 * time.Millisecond,
		wTO:     50 * time.Millisecond,
		clk:     timekit.ClockStartingAt(now),
		t:       t,
	}
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

func (tst *Tester) SrvCon() net.Conn  { return tst.srvCC }
func (tst *Tester) DataCon() net.Conn { return tst.cliCC }

func (tst *Tester) TstConfig(opts ...ftpsrv.Option) ftpsrv.Config {
	defaults := []ftpsrv.Option{
		ftpsrv.WithReadTimeout(tst.rTO),
		ftpsrv.WithWriteTimeout(tst.wTO),
		ftpsrv.WithClock(tst.clk),
	}
	return ftpsrv.NewConfig(append(defaults, opts...)...)
}

func (tst *Tester) TstSession(opts ...ftpsrv.Option) *ftpsrv.Session {
	ses := &ftpsrv.Session{
		ID:         uuid.Must(uuid.NewV7()),
		Cfg:        tst.TstConfig(opts...),
		LocalAddr:  tst.srvCC.LocalAddr(),
		RemoteAddr: tst.cliCC.RemoteAddr(),
		StartedAt:  tst.clk(),
	}
	return ses
}
