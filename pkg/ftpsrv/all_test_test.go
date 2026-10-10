package ftpsrv_test

import (
	"io"
	"strings"

	"github.com/ctx42/logkit/pkg/logkit"
	"github.com/ctx42/testing/pkg/assert"
	"github.com/ctx42/testing/pkg/tester"
)

// checkMsg is a convenience for checking a log message is equal to the string.
var checkMsg = logkit.CheckMsg

// CtrlMgs returns all logged control channel messages from the [logkit.Tester].
func CtrlMgs(tst *logkit.Tester) []string {
	var ret []string
	var filter = func(ent logkit.Entry) error {
		val, err := ent.Str("message")
		if err != nil {
			return err
		}
		if strings.HasPrefix(val, "> ") || strings.HasPrefix(val, "< ") {
			ret = append(ret, val)
			return nil
		}
		return logkit.ErrValue
	}
	tst.Filter(filter)
	return ret
}

// WaitForCose closes the c and waits for the log message up to 1s and returns.
func WaitForCose(t tester.T, c io.Closer, tst *logkit.Tester, msg string) {
	t.Helper()
	assert.NoError(t, c.Close())
	tst.WaitForAny("1s", checkMsg(msg))
}
