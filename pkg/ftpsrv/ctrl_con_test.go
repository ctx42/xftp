package ftpsrv

import (
	"fmt"
	"testing"

	"github.com/ctx42/xftp/pkg/ftpsrv/ftpsrvtest"
)

func Test_Name(t *testing.T) {
	// --- Given ---
	tlog, zlog := TstLogger(t)

	tst := ftpsrvtest.NewTester(t).WireUp()
	cfg := NewConfig()
	ses := NewSession(TstID(), cfg)

	cc := NewCtrlCon(ses, tst.SrvCon(), zlog).Listen()

	// --- When ---
	_ = cc

	// --- Then ---
	fmt.Println(tlog.String()) // TODO():
}
