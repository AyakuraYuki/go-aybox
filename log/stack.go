package log

import (
	"runtime"
	"sync"

	bytesPool "github.com/AyakuraYuki/go-aybox/log/bytes_pool"
)

var stacktracePool = sync.Pool{
	New: func() interface{} {
		return newProgramCounters(64)
	},
}

type programCounters struct {
	pcs []uintptr
}

func newProgramCounters(size int) *programCounters {
	return &programCounters{pcs: make([]uintptr, size)}
}

var bufferPool = bytesPool.NewPool()

func TakeStacktrace(skipOpt ...int) string {
	skip := 2
	if len(skipOpt) > 0 {
		skip = skipOpt[0]
	}

	buf := bufferPool.Get()
	defer buf.Free()

	pc := stacktracePool.Get().(*programCounters)
	defer stacktracePool.Put(pc)

	var numFrames int
	for {
		// Skip the call to runtime.Counters and takeStacktrace so that the
		// program counters start at the caller of takeStacktrace.
		numFrames = runtime.Callers(skip, pc.pcs)
		if numFrames < len(pc.pcs) {
			break
		}
		// Don't put the too-short counter slice back into the pool; this lets
		// the pool adjust if we consistently take deep stacktraces.
		pc = newProgramCounters(len(pc.pcs) * 2)
	}

	frames := runtime.CallersFrames(pc.pcs[:numFrames])

	// Note: On the last iteration, frames.Next() returns false, with a valid
	// frame, but we ignore this frame. The last frame is a runtime frame which
	// adds noise, since it's only either runtime.main or runtime.goexit.
	i := 0
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		if i != 0 {
			buf.AppendByte('\n')
		}
		i++
		buf.AppendString(frame.Function)
		buf.AppendByte('\n')
		buf.AppendByte('\t')
		buf.AppendString(frame.File)
		buf.AppendByte(':')
		buf.AppendInt(int64(frame.Line))
	}

	return buf.String()
}
