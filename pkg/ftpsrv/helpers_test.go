package ftpsrv

import (
	"testing"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_SplitCmdLine(t *testing.T) {
	tt := []struct {
		testN string

		line  string
		wCmd  string
		wArgs []string
	}{
		{"no args", "CMD", "CMD", nil},
		{"one arg", "CMD arg0", "CMD", []string{"arg0"}},
		{"two args", "CMD arg0 arg1", "CMD", []string{"arg0", "arg1"}},
		{"two args", "CMD  arg0  arg1", "CMD", []string{"arg0", "arg1"}},
		{"empty line", "", "", nil},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			hCmd, hArgs := SplitCmdLine(tc.line)

			// --- Then ---
			assert.Equal(t, tc.wCmd, hCmd)
			assert.Equal(t, tc.wArgs, hArgs)
		})
	}
}
