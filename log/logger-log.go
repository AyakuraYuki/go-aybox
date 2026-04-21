package log

import (
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

// Log carries per-call state for a single log entry.
// Obtain one via Logger.Debug / Logger.Info / Logger.Warn / Logger.Error.
//
// Every method is nil-safe: when a Logger's effective level excludes the
// requested level, the factory (e.g. Logger.Info) returns nil, and every
// chained call on nil *Log is a no-op — mirroring zerolog's nil *Event path
// so that disabled-level call sites cost essentially nothing.
type Log struct {
	event   *zerolog.Event
	depth   int
	level   zerolog.Level
	stack   bool
	traceId string
	logger  *Logger
}

// KV adds a string key-value pair to this log entry.
func (b *Log) KV(key, val string) *Log {
	if b == nil {
		return nil
	}
	b.event = b.event.Str(key, val)
	return b
}

// TraceID prepends a trace ID to the message.
func (b *Log) TraceID(traceId string) *Log {
	if b == nil {
		return nil
	}
	b.traceId = traceId
	return b
}

// Stack enables a stack trace field ("stack") on this log entry.
func (b *Log) Stack() *Log {
	if b == nil {
		return nil
	}
	b.stack = true
	return b
}

// Err attaches an error to the log entry.
func (b *Log) Err(err error) *Log {
	if b == nil {
		return nil
	}
	b.event = b.event.Err(err)
	return b
}

// Bool adds a bool key-value pair to this log entry.
func (b *Log) Bool(key string, val bool) *Log {
	if b == nil {
		return nil
	}
	b.event = b.event.Bool(key, val)
	return b
}

// Int adds an int key-value pair to this log entry.
func (b *Log) Int(key string, val int) *Log {
	if b == nil {
		return nil
	}
	b.event = b.event.Int(key, val)
	return b
}

// Int32 adds an int32 key-value pair to this log entry.
func (b *Log) Int32(key string, val int32) *Log {
	if b == nil {
		return nil
	}
	b.event = b.event.Int32(key, val)
	return b
}

// Int64 adds an int64 key-value pair to this log entry.
func (b *Log) Int64(key string, val int64) *Log {
	if b == nil {
		return nil
	}
	b.event = b.event.Int64(key, val)
	return b
}

// Uint adds a uint key-value pair to this log entry.
func (b *Log) Uint(key string, val uint) *Log {
	if b == nil {
		return nil
	}
	b.event = b.event.Uint(key, val)
	return b
}

// Uint32 adds a uint32 key-value pair to this log entry.
func (b *Log) Uint32(key string, val uint32) *Log {
	if b == nil {
		return nil
	}
	b.event = b.event.Uint32(key, val)
	return b
}

// Uint64 adds a uint64 key-value pair to this log entry.
func (b *Log) Uint64(key string, val uint64) *Log {
	if b == nil {
		return nil
	}
	b.event = b.event.Uint64(key, val)
	return b
}

// Float32 adds a float32 key-value pair to this log entry.
func (b *Log) Float32(key string, val float32) *Log {
	if b == nil {
		return nil
	}
	b.event = b.event.Float32(key, val)
	return b
}

// Float64 adds a float64 key-value pair to this log entry.
func (b *Log) Float64(key string, val float64) *Log {
	if b == nil {
		return nil
	}
	b.event = b.event.Float64(key, val)
	return b
}

// Str adds a string key-value pair to this log entry.
func (b *Log) Str(key, val string) *Log {
	if b == nil {
		return nil
	}
	b.event = b.event.Str(key, val)
	return b
}

// Strs adds a []string key-value pair to this log entry.
func (b *Log) Strs(key string, vals []string) *Log {
	if b == nil {
		return nil
	}
	b.event = b.event.Strs(key, vals)
	return b
}

// Interface adds an any key-value pair to this log entry.
func (b *Log) Interface(key string, val any) *Log {
	if b == nil {
		return nil
	}
	b.event = b.event.Interface(key, val)
	return b
}

// Time adds a time.Time key-value pair to this log entry.
func (b *Log) Time(key string, val time.Time) *Log {
	if b == nil {
		return nil
	}
	b.event = b.event.Time(key, val)
	return b
}

// Dur adds a time.Duration key-value pair to this log entry.
func (b *Log) Dur(key string, val time.Duration) *Log {
	if b == nil {
		return nil
	}
	b.event = b.event.Dur(key, val)
	return b
}

// Array adds the field key with an array to the event context.
// Use zerolog.Arr() to create the array or pass a type that
// implement the LogArrayMarshaler interface.
func (b *Log) Array(key string, arr zerolog.LogArrayMarshaler) *Log {
	if b == nil {
		return nil
	}
	b.event = b.event.Array(key, arr)
	return b
}

// Event returns the underlying *zerolog.Event for direct zerolog API access.
// Calling .Msg() on the returned event bypasses TraceID, codeline, and stack
// trace. Prefer Log.Msg or Log.Msgf to emit with full instrumentation.
func (b *Log) Event() *zerolog.Event {
	if b == nil {
		return nil
	}
	return b.event
}

// Msg emits the log entry with a plain string message.
// This is the zero-alloc terminal method: it uses event.Msg(string) internally,
// avoiding any fmt.Sprintf overhead.
//
// The *Log is returned to the pool after emission; do not use b after Msg returns.
func (b *Log) Msg(msg string) {
	if b == nil {
		return
	}
	if b.traceId != "" {
		msg = "[trace_id: " + b.traceId + "] " + msg
	}
	if b.logger.codeline && b.depth != 0 {
		lineDesc, fn := codeline(b.depth)
		b.event = b.event.Str("codeline", lineDesc).Str("func", fn)
	}
	if b.stack {
		b.event = b.event.Str("stack", TakeStacktrace(b.depth+1))
	}
	b.event.Msg(msg)
	b.recycle()
}

// Msgf emits the log entry with a formatted message.
// When called with no format arguments, it takes the same zero-alloc path as
// Msg. The fmt.Sprintf cost is only paid when format arguments are present.
//
// The *Log is returned to the pool after emission; do not use b after Msgf returns.
func (b *Log) Msgf(msg string, v ...any) {
	if b == nil {
		return
	}
	if b.traceId != "" {
		msg = "[trace_id: " + b.traceId + "] " + msg
	}
	if b.logger.codeline && b.depth != 0 {
		lineDesc, fn := codeline(b.depth)
		b.event = b.event.Str("codeline", lineDesc).Str("func", fn)
	}
	if b.stack {
		b.event = b.event.Str("stack", TakeStacktrace(b.depth+1))
	}
	if len(v) == 0 {
		b.event.Msg(msg)
	} else {
		b.event.Msgf(msg, v...)
	}
	b.recycle()
}

// recycle clears heap-retaining fields and returns b to logPool.
// Must only be called once per *Log, at the end of Msg or Msgf.
func (b *Log) recycle() {
	b.event = nil  // release zerolog event reference
	b.traceId = "" // release string reference
	b.logger = nil // release logger reference
	logPool.Put(b)
}

// codeline returns the source file:line and function name at call depth n,
// using strconv to avoid fmt.Sprintf overhead.
func codeline(n int) (lineDesc string, fn string) {
	pc, file, line, ok := runtime.Caller(n)
	if !ok {
		return "", ""
	}
	if i := strings.Index(file, "src/"); i >= 0 {
		file = file[i+4:]
	}
	buf := make([]byte, 0, len(file)+12)
	buf = append(buf, file...)
	buf = append(buf, ':')
	buf = strconv.AppendInt(buf, int64(line), 10)
	return string(buf), runtime.FuncForPC(pc).Name()
}
