package ftpsrv_test

import (
	"testing"

	"github.com/ctx42/testing/pkg/assert"

	"github.com/ctx42/xftp/pkg/ftpsrv/ftpsrvtest"

	. "github.com/ctx42/xftp/pkg/ftpsrv"
)

func Test_CtrlCon_Listen(t *testing.T) {
	t.Run("starts and sends a server ready message", func(t *testing.T) {
		// --- Given ---
		tst := ftpsrvtest.NewTester(t).WireUp()
		ses := tst.TstSession()
		cc := NewCtrlCon(ses, tst.SrvCon(), tst.Logger()).Listen()
		t.Cleanup(func() { t.Helper(); assert.NoError(t, cc.Close()) })

		// --- When ---
		cc.Listen()

		// --- Then ---
		tst.Reply(ServerReady.String())

		tlog := tst.ExamineLog()
		tlog.WaitForAny("1s", hasMsg("cc.listen: started"))
		tlog.WaitForAny("1s", hasMsg(ServerReady.String()))
	})
}
