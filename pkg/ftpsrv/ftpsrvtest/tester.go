package ftpsrvtest

import (
	"net"
	"time"

	"github.com/ctx42/testing/pkg/tester"
)

type Tester struct {
	cliCC   net.Conn      // Client side of the control connection.
	cliDC   net.Conn      // Client side of a data connection.
	srvCC   net.Conn      // Server side of the control connection.
	replays []string      // Lines read from the client side control connection.
	rTO     time.Duration // Timeout reading from control connection.
	wTO     time.Duration // Timeout writing to control connection.
	t       tester.T
}

func NewTester(t tester.T) *Tester {
	t.Helper()
	tst := &Tester{
		replays: make([]string, 0, 10),
		rTO:     50 * time.Millisecond,
		wTO:     50 * time.Millisecond,
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

func (tst *Tester) SrvConn() net.Conn { return tst.srvCC }
