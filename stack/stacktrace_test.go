package stack

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func a() *Stacktrace { return b() }
func b() *Stacktrace { return c() }
func c() *Stacktrace { return d() }
func d() *Stacktrace { return e() }
func e() *Stacktrace { return f() }
func f() *Stacktrace { return NewStacktrace("fff") }

func TestStacktrace(t *testing.T) {
	st := a()
	assert.NotNil(t, st)
	assert.Equal(t, "fff", st.span)

	fmt.Println(st.String(""))
	fmt.Println()
	fmt.Println(st.Sources())
}
