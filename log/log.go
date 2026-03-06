package log

import (
	"fmt"
	"io"
	"runtime"
	"strings"
	"time"

	"code.cloudfoundry.org/go-diodes"
	"github.com/rs/zerolog"

	"github.com/AyakuraYuki/go-aybox/ip"
	"github.com/AyakuraYuki/go-aybox/log/async"
	"github.com/AyakuraYuki/go-aybox/log/console"
)

func init() {
	zerolog.MessageFieldName = "desc"
	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimeFieldFormat = ""
}

// CloseFn is a function that closes a resource.
type CloseFn func() error

// Leveler is implemented by writers that expose a minimum log level.
type Leveler interface {
	Level() zerolog.Level
}

// Logger is a multi-instance logger backed by zerolog.
// Create one with New(); safe to call methods from multiple goroutines.
type Logger struct {
	zlogger  zerolog.Logger
	depth    int  // base call-stack depth for codeline; default 2
	codeline bool // emit file:line and func fields; default false
	closeFns []CloseFn
}

// New creates a Logger with the provided options.
// If no WithWriters option is given, a ConsoleWriter is used.
func New(opts ...Option) *Logger {
	cfg := &config{
		level:    zerolog.DebugLevel,
		depth:    2,
		hostname: ip.Hostname(),
	}
	for _, opt := range opts {
		opt(cfg)
	}

	writers := cfg.writers
	if len(writers) == 0 {
		writers = []io.Writer{console.New()}
	}

	var closeFns []CloseFn
	if cfg.async {
		wrapped := make([]io.Writer, len(writers))
		for i, writer := range writers {
			lvl := cfg.level
			if lw, ok := writer.(Leveler); ok {
				lvl = lw.Level()
			}
			aw := async.New(lvl, writer,
				diodes.NewManyToOne(1024, diodes.AlertFunc(func(missed int) {
					fmt.Printf("logger dropped %d messages\n", missed)
				})),
				time.Second)
			wrapped[i] = aw
			closeFns = append(closeFns, aw.Close)
		}
		writers = wrapped
	}

	var writer io.Writer
	if len(writers) == 1 {
		writer = writers[0]
	} else {
		writer = zerolog.MultiLevelWriter(writers...)
	}

	logger := zerolog.New(writer).With().
		Timestamp().
		Str("host", cfg.hostname).
		Logger()

	if len(cfg.fields) > 0 {
		ctx := logger.With()
		for k, v := range cfg.fields {
			ctx = ctx.Str(k, v)
		}
		logger = ctx.Logger()
	}

	return &Logger{
		zlogger:  logger,
		depth:    cfg.depth,
		codeline: cfg.codeline,
		closeFns: closeFns,
	}
}

// Close flushes and closes all writers registered to this logger.
func (l *Logger) Close() {
	for _, fn := range l.closeFns {
		_ = fn()
	}
}

// Log starts a new message with no level.
func (l *Logger) Log(name ...string) *Log { return l.newLog(zerolog.NoLevel, name...) }

// Debug returns a *Log at debug level, optionally tagged with name.
func (l *Logger) Debug(name ...string) *Log { return l.newLog(zerolog.DebugLevel, name...) }

// Info returns a *Log at info level, optionally tagged with name.
func (l *Logger) Info(name ...string) *Log { return l.newLog(zerolog.InfoLevel, name...) }

// Warn returns a *Log at warn level, optionally tagged with name.
func (l *Logger) Warn(name ...string) *Log { return l.newLog(zerolog.WarnLevel, name...) }

// Error returns a *Log at error level, optionally tagged with name.
func (l *Logger) Error(name ...string) *Log { return l.newLog(zerolog.ErrorLevel, name...) }

// newLog creates a *Log for the given level.
// Returns nil without allocating when the level is disabled — same fast path
// that zerolog uses for a nil *Event.
func (l *Logger) newLog(level zerolog.Level, name ...string) *Log {
	event := l.zlogger.WithLevel(level)
	if event == nil {
		return nil
	}
	if len(name) > 0 && name[0] != "" {
		event = event.Str("name", name[0])
	}
	return &Log{
		event:  event,
		depth:  l.depth,
		level:  level,
		logger: l,
	}
}

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
func (b *Log) KV(key, val string) *Log { //nolint:unparam
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

// Err attaches an error field (only effective at ErrorLevel).
func (b *Log) Err(err error) *Log {
	if b == nil {
		return nil
	}
	if b.level == zerolog.ErrorLevel {
		b.event = b.event.Err(err)
	}
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

// Event returns the underlying *zerolog.Event for direct zerolog API access.
// Calling .Msg() on the returned event bypasses TraceID, codeline, and stack
// trace. Prefer Log.Msg or Log.Msgf to emit with full instrumentation.
func (b *Log) Event() *zerolog.Event {
	if b == nil {
		return nil
	}
	return b.event
}

// Msg outputs a log message from one or more values.
func (b *Log) Msg(msg string) {
	if b == nil {
		return
	}
	b.depth++
	b.Msgf(msg)
}

// Msgf outputs a formatted log message.
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
	b.event.Msgf(msg, v...)
}

func codeline(n int) (lineDesc string, fn string) {
	pc, file, line, ok := runtime.Caller(n)
	if !ok {
		return "", ""
	}
	if i := strings.Index(file, "src/"); i >= 0 {
		return fmt.Sprintf("%s:%d", file[i+4:], line), runtime.FuncForPC(pc).Name()
	}
	return fmt.Sprintf("%s:%d", file, line), runtime.FuncForPC(pc).Name()
}
