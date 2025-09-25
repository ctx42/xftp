package ftpsrv_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ctx42/xftp/pkg/ftpsrv/ftpsrvtest"

	. "github.com/ctx42/xftp/pkg/ftpsrv"
)

func Test_Name(t *testing.T) {
	// --- Given ---
	tlog, zlog := TstLogger(t)

	tst := ftpsrvtest.NewTester(t).WireUp()
	ses := tst.TstSession()

	cc := NewCtrlCon(ses, tst.SrvCon(), zlog).Listen()

	// --- When ---
	_ = cc

	// --- Then ---
	time.Sleep(1 * time.Second)
	fmt.Println(tlog.String()) // TODO():
}
