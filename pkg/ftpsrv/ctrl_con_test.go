package ftpsrv_test

import (
	"testing"
	"time"

	"github.com/ctx42/testing/pkg/assert"

	"github.com/ctx42/xftp/pkg/ftpcmd"
	"github.com/ctx42/xftp/pkg/ftpsrv/ftpsrvtest"

	. "github.com/ctx42/xftp/pkg/ftpsrv"
)

func Test_NewCtrlCon(t *testing.T) {
	// --- Given ---
	tst := ftpsrvtest.NewTester(t).WireUp()
	ses := tst.Session()

	// --- When ---
	cc := NewCtrlCon(ses, tst.SrvCon(), tst.Logger())

	// --- Then ---
	assert.NoError(t, cc.Close())
	tst.ExamineLog().Entries().AssertLen(0)
}

func Test_CtrlCon_Listen(t *testing.T) {
	t.Run("starts and sends a server ready message", func(t *testing.T) {
		// --- Given ---
		tst := ftpsrvtest.NewTester(t).WireUp()
		ses := tst.Session()

		cc := NewCtrlCon(ses, tst.SrvCon(), tst.Logger())
		tst.CloseAfterTest(cc)

		// --- When ---
		cc.Listen()

		// --- Then ---
		tst.GetReply(ServerReady.String())

		tlog := tst.ExamineLog()
		WaitForCose(t, cc, tlog, "cc.listen: exiting")

		hCCL := CtrlMgs(tlog)
		wCCL := []string{"> 220 FTP Server ready."}
		assert.Equal(t, hCCL, wCCL)
	})

	t.Run("handles command", func(t *testing.T) {
		// --- Given ---
		tst := ftpsrvtest.NewTester(t).WireUp()
		ses := tst.Session()

		cc := NewCtrlCon(ses, tst.SrvCon(), tst.Logger())
		tst.CloseAfterTest(cc)

		// --- When ---
		cc.Listen()

		// --- Then ---
		tst.GetReply(ServerReady.String())
		tst.SendCmd(ftpcmd.NOOP, "").GetReply(NOOPSuccess.String())

		tlog := tst.ExamineLog()
		WaitForCose(t, cc, tlog, "cc.listen: exiting")

		hCCL := CtrlMgs(tlog)
		wCCL := []string{
			"> 220 FTP Server ready.",
			"< NOOP",
			"> 200 NOOP OK.",
		}
		assert.Equal(t, hCCL, wCCL)
	})

	t.Run("client closes the connection after connecting", func(t *testing.T) {
		// --- Given ---
		tst := ftpsrvtest.NewTester(t).WireUp()
		ses := tst.Session()

		cc := NewCtrlCon(ses, tst.SrvCon(), tst.Logger())
		tst.CloseAfterTest(cc)

		// --- When ---
		cc.Listen()

		// --- Then ---
		tst.CliConClose()

		time.Sleep(time.Second)
		tlog := tst.ExamineLog()
		want := "" +
			"cc.write_line fail: " +
			"io: read/write on closed pipe; " +
			"msg: 220 FTP Server ready."
		tlog.WaitForAny("1s", checkMsg(want))
		tlog.WaitForAny("1s", checkMsg("cc.listen: exiting"))
		assert.Equal(t, []string{}, CtrlMgs(tlog))
	})

	t.Run("client closes the connection after the welcome message", func(t *testing.T) {
		// --- Given ---
		tst := ftpsrvtest.NewTester(t).WireUp()
		ses := tst.Session()

		cc := NewCtrlCon(ses, tst.SrvCon(), tst.Logger())
		tst.CloseAfterTest(cc)

		// --- When ---
		cc.Listen()

		// --- Then ---
		tst.GetReply(ServerReady.String())
		tst.CliConClose()

		time.Sleep(time.Second)
		tlog := tst.ExamineLog()
		want := "cc.listen: connection closed by client"
		tlog.WaitForAny("1s", checkMsg(want))
		tlog.WaitForAny("1s", checkMsg("cc.listen: exiting"))
		assert.Equal(t, []string{"> 220 FTP Server ready."}, CtrlMgs(tlog))
	})
}
