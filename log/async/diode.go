package async

import (
	"context"
	"io"
	"time"

	"code.cloudfoundry.org/go-diodes"
	"github.com/rs/zerolog"

	bytesPool "github.com/AyakuraYuki/go-aybox/log/bytes_pool"
)

var diodeBufPool = bytesPool.NewPool()

// Writer is a io.Writer wrapper that uses a diode to make Write lock-free,
// non-blocking and thread safe.
type Writer struct {
	l            zerolog.Level
	w            io.Writer
	d            *diodes.ManyToOne
	p            *diodes.Poller
	c            context.CancelFunc
	done         chan struct{}
	closeTimeout time.Duration
}

// New creates a writer wrapping w with a many-to-one diode in order
// to never block log producers and drop events if the writer can't keep up
// with the flow of data.
//
// closeTimeout overrides the default drain timeout used by Close. When set to
// a positive value it replaces the 30-second default; zero or negative leaves
// the default in effect.
//
// See code.cloudfoundry.org/go-diodes for more info on diode.
func New(
	l zerolog.Level,
	w io.Writer,
	manyToOneDiode *diodes.ManyToOne,
	poolInterval time.Duration,
	closeTimeout time.Duration,
) Writer {
	ctx, cancel := context.WithCancel(context.Background())
	dw := Writer{
		l: l,
		w: w,
		d: manyToOneDiode,
		p: diodes.NewPoller(manyToOneDiode,
			diodes.WithPollingInterval(poolInterval),
			diodes.WithPollingContext(ctx)),
		c:            cancel,
		done:         make(chan struct{}),
		closeTimeout: closeTimeout,
	}
	go dw.poll()
	return dw
}

// Write copies p into a pooled Buffer and enqueues it for the background writer.
func (dw Writer) Write(p []byte) (n int, err error) {
	buf := diodeBufPool.Get()
	_, _ = buf.Write(p)
	dw.d.Set(diodes.GenericDataType(buf))
	return len(p), nil
}

// WriteLevel writes data to writer with level info provided.
func (dw Writer) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	if level < dw.l {
		return len(p), nil
	}
	return dw.Write(p)
}

// Close signals the background poll goroutine to stop and waits for it to
// drain any buffered entries. The wait is bounded by closeTimeout when it is
// positive, or by the 30-second default otherwise; if the goroutine has not
// exited by then, Close returns early and any remaining buffered entries are
// lost. The underlying writer is closed afterwards if it implements io.Closer.
func (dw Writer) Close() error {
	dw.c()
	closeTimeout := 30 * time.Second
	if dw.closeTimeout > 0 {
		closeTimeout = dw.closeTimeout
	}
	select {
	case <-dw.done:
	case <-time.After(closeTimeout):
	}
	if w, ok := dw.w.(io.Closer); ok {
		return w.Close()
	}
	return nil
}

func (dw Writer) poll() {
	defer close(dw.done)
	for {
		d := dw.p.Next()
		if d == nil {
			return
		}
		buf := (*bytesPool.Buffer)(d)
		_, _ = dw.w.Write(buf.Bytes())
		buf.Free()
	}
}
