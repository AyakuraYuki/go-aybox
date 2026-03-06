package log

import (
	"context"
	"io"
	"time"

	"code.cloudfoundry.org/go-diodes"
	"github.com/rs/zerolog"
)

var diodeBufPool = NewBytesPool()

// Writer is a io.Writer wrapper that uses a diode to make Write lock-free,
// non-blocking and thread safe.
type Writer struct {
	l    zerolog.Level
	w    io.Writer
	d    *diodes.ManyToOne
	p    *diodes.Poller
	c    context.CancelFunc
	done chan struct{}
}

// NewAsyncWriter creates a writer wrapping w with a many-to-one diode in order
// to never block log producers and drop events if the writer can't keep up
// with the flow of data.
//
// Use a diode.Writer when
//
//	d := diodes.NewManyToOne(1000, diodes.AlertFunc(func(missed int) {
//	    log.Printf("Dropped %d messages", missed)
//	}))
//	w := diode.NewWriter(w, d, 10 * time.Millisecond)
//	log := zerolog.New(w)
//
// See code.cloudfoundry.org/go-diodes for more info on diode.
func NewAsyncWriter(
	l zerolog.Level,
	w io.Writer,
	manyToOneDiode *diodes.ManyToOne,
	poolInterval time.Duration,
) Writer {
	ctx, cancel := context.WithCancel(context.Background())
	dw := Writer{
		l: l,
		w: w,
		d: manyToOneDiode,
		p: diodes.NewPoller(manyToOneDiode,
			diodes.WithPollingInterval(poolInterval),
			diodes.WithPollingContext(ctx)),
		c:    cancel,
		done: make(chan struct{}),
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

// Close releases the diode poller and calls Close on the wrapped writer if
// io.Closer is implemented.
func (dw Writer) Close() error {
	dw.c()
	<-dw.done
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
		buf := (*Buffer)(d)
		_, _ = dw.w.Write(buf.Bytes())
		buf.Free()
	}
}
