package lines

import (
	"runtime"
	"testing"
)

func TestLineSep(t *testing.T) {
	t.Logf("os: %s, line sep: %#v", runtime.GOOS, LineSep)
}
