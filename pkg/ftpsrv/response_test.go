package ftpsrv

import (
	"testing"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_Resp(t *testing.T) {
	// --- When ---
	rsp := Resp(200, "some", "lines")

	// --- Then ---
	assert.Equal(t, 200, rsp.code)
	assert.Equal(t, []string{"some", "lines"}, rsp.lines)
}

func Test_Response_String(t *testing.T) {
	t.Run("single line", func(t *testing.T) {
		// --- Given ---
		rsp := Resp(200, "line0")

		// --- When ---
		have := rsp.String()

		// --- Then ---
		assert.Equal(t, "200 line0", have)
	})

	t.Run("two lines", func(t *testing.T) {
		// --- Given ---
		rsp := Resp(200, "line0", "line1")

		// --- When ---
		have := rsp.String()

		// --- Then ---
		want := "200-line0\r\n200 line1"
		assert.Equal(t, want, have)
	})

	t.Run("three lines", func(t *testing.T) {
		// --- Given ---
		rsp := Resp(200, "line0", "line1", "line2")

		// --- When ---
		have := rsp.String()

		// --- Then ---
		want := "200-line0\r\nline1\r\n200 line2"
		assert.Equal(t, want, have)
	})

	t.Run("three lines middle one starts with numbers", func(t *testing.T) {
		// --- Given ---
		rsp := Resp(200, "line0", "123 line1", "line2")

		// --- When ---
		have := rsp.String()

		// --- Then ---
		want := "200-line0\r\n 123 line1\r\n200 line2"
		assert.Equal(t, want, have)
	})

	t.Run("three lines last one starts with numbers", func(t *testing.T) {
		// --- Given ---
		rsp := Resp(200, "line0", "line1", "123 line2")

		// --- When ---
		have := rsp.String()

		// --- Then ---
		want := "200-line0\r\nline1\r\n200 123 line2"
		assert.Equal(t, want, have)
	})

	t.Run("pad middle lines", func(t *testing.T) {
		// --- Given ---
		rsp := Resp(200, "line0", "line1", "line2").Pad()

		// --- When ---
		have := rsp.String()

		// --- Then ---
		want := "200-line0\r\n line1\r\n200 line2"
		assert.Equal(t, want, have)
	})

	t.Run("no double pad lines starting with number", func(t *testing.T) {
		// --- Given ---
		rsp := Resp(200, "line0", "123 line1", "line2").Pad()

		// --- When ---
		have := rsp.String()

		// --- Then ---
		want := "200-line0\r\n 123 line1\r\n200 line2"
		assert.Equal(t, want, have)
	})

	t.Run("do not add CRLF twice", func(t *testing.T) {
		// --- Given ---
		rsp := Resp(200, "line0\r\n", "line1\n", "line2\r\n")

		// --- When ---
		have := rsp.String()

		// --- Then ---
		want := "200-line0\r\nline1\r\n200 line2"
		assert.Equal(t, want, have)
	})

	t.Run("not initialized instance", func(t *testing.T) {
		// --- Given ---
		rsp := &Response{}

		// --- When ---
		have := rsp.String()

		// --- Then ---
		assert.Empty(t, have)
	})
}

func Test_Response_With(t *testing.T) {
	t.Run("three lines last one starts with numbers", func(t *testing.T) {
		// --- Given ---
		rsp := Resp(200, "line%d", "line%d", "123 line%d")

		// --- When ---
		have := rsp.With(0, 1, 2)

		// --- Then ---
		want := "200-line0\r\nline1\r\n200 123 line2"
		assert.Equal(t, want, have)
	})
}

func Test_Response_Text(t *testing.T) {
	t.Run("single line", func(t *testing.T) {
		// --- Given ---
		rsp := Resp(200, "line0")

		// --- When ---
		have := rsp.Text()

		// --- Then ---
		want := "200 line0"
		assert.Equal(t, want, have)
	})

	t.Run("two lines", func(t *testing.T) {
		// --- Given ---
		rsp := Resp(200, "line0", "line1")

		// --- When ---
		have := rsp.Text()

		// --- Then ---
		want := "200 line0\nline1"
		assert.Equal(t, want, have)
	})

	t.Run("three lines", func(t *testing.T) {
		// --- Given ---
		rsp := Resp(200, "line0", "line1", "123 line2")

		// --- When ---
		have := rsp.Text()

		// --- Then ---
		want := "200 line0\nline1\n123 line2"
		assert.Equal(t, want, have)
	})
}
