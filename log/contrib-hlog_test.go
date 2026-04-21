package log

import (
	"math/rand"
	"runtime"
	"testing"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"

	"github.com/AyakuraYuki/go-aybox/randoms"
)

func TestHLog(t *testing.T) {
	hlog.SetLogger(defaultLogger)

	hlog.Tracef("trace: %v", randoms.RandPassword(20))
	hlog.Trace(rand.Int63n(100))

	hlog.Debugf("now: %v", time.Now().Format(time.DateTime))
	hlog.Debug("hey!")

	hlog.Infof("hi, %s", randoms.RandString(10, randoms.LowerUpper))
	hlog.Info("hello, world")

	hlog.Noticef("dice: %d", rand.Int63n(100))
	hlog.Notice("ping-pong")

	hlog.Warnf("FXI: %v", []string{"open", "the", "door"})
	hlog.Warn("joking")

	hlog.Errorf("oops: %s", "cake is a lie")
	hlog.Error(runtime.NumCgoCall())

	hlog.Fatalf("panic: %d", rand.Int63n(100))
	hlog.Fatal("abc is not xyz")
}
