package stack

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func a() *Stacktrace { return b() }
func b() *Stacktrace { return c() }
func c() *Stacktrace { return d() }
func d() *Stacktrace { return e() }
func e() *Stacktrace { return f() }
func f() *Stacktrace {
	err := errors.New("test error")
	return NewStacktrace(err, "fff")
}

func TestStacktrace(t *testing.T) {
	st := a()
	assert.NotNil(t, st)
	assert.Equal(t, "fff", st.span)
}
