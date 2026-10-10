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
		tst := ftpsrvtest.NewTester(t).WireUp().INE()
		ses := tst.Session()
		cfg := tst.Config()

		cc := NewCtrlCon(ses, cfg, tst.SrvCon(), tst.Logger()).Listen()
		tst.CloseAfterTest(cc)
		tst.GetReply(ServerReady.String())

		// --- When ---
		tst.SendCmd(ftpcmd.NOOP, "")

		// --- Then ---
		tst.GetReply(NOOPSuccess.String())
	})

	t.Run("wrong number of arguments", func(t *testing.T) {
		// --- Given ---
		tst := ftpsrvtest.NewTester(t).WireUp().INE()
		ses := tst.Session()
		cfg := tst.Config()

		cc := NewCtrlCon(ses, cfg, tst.SrvCon(), tst.Logger()).Listen()
		tst.CloseAfterTest(cc)
		tst.GetReply(ServerReady.String())

		// --- When ---
		tst.SendCmd(ftpcmd.NOOP, "abc")

		// --- Then ---
		tst.GetReply(ErrorArgNum.String())
	})
}
