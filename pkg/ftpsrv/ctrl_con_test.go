package ftpsrv

import (
	"testing"

	"github.com/ctx42/xftp/pkg/ftpsrv/ftpsrvtest"
)

func Test_Name(t *testing.T) {
	// --- Given ---
	tst := ftpsrvtest.NewTester(t).WireUp()
	cfg := NewConfig()

	cc := NewCtrlCon(cfg, tst.SrvConn()).Listen()

	// --- When ---
	_ = cc

	// --- Then ---
}
