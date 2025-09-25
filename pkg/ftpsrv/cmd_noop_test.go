package ftpsrv_test

import (
	"testing"

	"github.com/ctx42/xftp/pkg/ftpcmd"
	"github.com/ctx42/xftp/pkg/ftpsrv/ftpsrvtest"

	. "github.com/ctx42/xftp/pkg/ftpsrv"
)

func Test_handleNOOP(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		tst := ftpsrvtest.NewTester(t).WireUp()
		ses := tst.TstSession()

		cc := NewCtrlCon(ses, tst.SrvCon(), tst.Logger()).Listen()
		tst.CloseAfterTest(cc)
		tst.Reply(ServerReady.String())

		// --- When ---
		tst.Cmd(ftpcmd.NOOP, "")

		// --- Then ---
		tst.Reply(NOOPSuccess.String())
	})

	t.Run("wrong number of arguments", func(t *testing.T) {
		// --- Given ---
		tst := ftpsrvtest.NewTester(t).WireUp()
		ses := tst.TstSession()

		cc := NewCtrlCon(ses, tst.SrvCon(), tst.Logger()).Listen()
		tst.CloseAfterTest(cc)
		tst.Reply(ServerReady.String())

		// --- When ---
		tst.Cmd(ftpcmd.NOOP, "abc")

		// --- Then ---
		tst.Reply(ErrorArgNum.String())
	})
}
