package log

import (
	"context"
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

type Leveler interface {
	Level() zerolog.Level
}

// Logger is a multi-instance logger backed by zerolog.
// Create one with New(); never share a pointer across goroutines without synchronisation.
type Logger struct {
	zlogger  zerolog.Logger
	depth    int
	closeFns []CloseFn
}

// New creates a Logger with the provided options.
// If no WithOutput option is given, a ConsoleWriter is used.
func New(opts ...Option) *Logger {
	cfg := &config{
		depth:    2,
		hostname: ip.Hostname(),
		level:    zerolog.DebugLevel,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	writers := cfg.writers
	if len(writers) == 0 {
		// use console writer as default
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

	var out io.Writer
	if len(writers) == 1 {
		out = writers[0]
	} else {
		out = zerolog.MultiLevelWriter(writers...)
	}

	zl := zerolog.New(out).With().Timestamp().Str("host", cfg.hostname).Logger()
	if len(cfg.fields) > 0 {
		c := zl.With()
		for k, v := range cfg.fields {
			c = c.Str(k, v)
		}
		zl = c.Logger()
	}

	return &Logger{
		zlogger:  zl,
		depth:    cfg.depth,
		closeFns: closeFns,
	}
}

// Close flushes and closes all writers registered to this logger.
func (l *Logger) Close() {
	for _, fn := range l.closeFns {
		_ = fn()
	}
}

// Debug returns a *Log at debug level, optionally tagged with name.
func (l *Logger) Debug(name ...string) *Log {
	return l.newLog(l.depth, zerolog.DebugLevel, name...)
}

// Info returns a *Log at info level, optionally tagged with name.
func (l *Logger) Info(name ...string) *Log {
	return l.newLog(l.depth, zerolog.InfoLevel, name...)
}

// Warn returns a *Log at warn level, optionally tagged with name.
func (l *Logger) Warn(name ...string) *Log {
	return l.newLog(l.depth, zerolog.WarnLevel, name...)
}

// Error returns a *Log at error level, optionally tagged with name.
func (l *Logger) Error(name ...string) *Log {
	return l.newLog(l.depth, zerolog.ErrorLevel, name...)
}

// Sample returns a *LogContext with the given sampler applied to every entry.
func (l *Logger) Sample(sampler zerolog.Sampler) *Context {
	return &Context{logger: l, sampler: sampler}
}

// FromContext extracts a *LogContext stored by (*LogContext).ToContext.
// Falls back to a plain *LogContext backed by this logger if none is found.
func (l *Logger) FromContext(ctx context.Context) *Context {
	if lc, ok := ctx.Value(contextKey{}).(*Context); ok {
		return lc
	}
	return &Context{logger: l}
}

func (l *Logger) newLog(depth int, level zerolog.Level, name ...string) *Log {
	zl := l.zlogger
	if len(name) > 0 && name[0] != "" {
		zl = zl.With().Str("name", name[0]).Logger()
	}
	return &Log{
		depth:   depth,
		level:   level,
		zlogger: &zl,
		logger:  l,
	}
}

// Log carries per-call state for a single log entry.
// Obtain one via Logger.Debug / Logger.Info / Logger.Warn / Logger.Error.
type Log struct {
	depth   int
	level   zerolog.Level
	stack   bool
	err     error
	traceId string
	zlogger *zerolog.Logger
	sampler zerolog.Sampler
	logger  *Logger
}

// Context stores this Log's zerolog state into a Go context (with context.Background as parent).
// Retrieve it later with (*Logger).FromContext.
func (b *Log) Context() context.Context {
	lc := &Context{logger: b.logger, zlogger: b.zlogger, sampler: b.sampler}
	return context.WithValue(context.Background(), contextKey{}, lc)
}

// KV adds a string key-value pair to this log entry.
func (b *Log) KV(key string, val string) *Log {
	zl := b.zlogger.With().Str(key, val).Logger()
	b.zlogger = &zl
	return b
}

// TraceID prepends a trace ID to the message.
func (b *Log) TraceID(traceId string) *Log {
	b.traceId = traceId
	return b
}

// Stack enables stack trace appended to the message.
func (b *Log) Stack() *Log {
	b.stack = true
	return b
}

// Err attaches an error field (only effective at ErrorLevel).
func (b *Log) Err(err error) *Log {
	if b.level == zerolog.ErrorLevel {
		b.err = err
	}
	return b
}

// Event returns a *zerolog.Event for direct zerolog API access.
func (b *Log) Event() *zerolog.Event {
	zl := *b.zlogger
	if b.sampler != nil {
		zl = zl.Sample(b.sampler)
	}
	return zl.WithLevel(b.level)
}

// Msg outputs a log message from one or more values.
func (b *Log) Msg(msg ...any) {
	b.depth++
	switch len(msg) {
	case 0:
		return
	case 1:
		switch v := msg[0].(type) {
		case string:
			b.Msgf(v)
		default:
			b.Msgf("%v", v)
		}
	default:
		fmtStr := strings.Repeat("%v, ", len(msg))
		b.Msgf(fmtStr[:len(fmtStr)-2], msg...)
	}
}

// Msgf outputs a formatted log message.
func (b *Log) Msgf(msg string, v ...any) {
	if b.traceId != "" {
		msg = fmt.Sprintf("traceId:[%s] ", b.traceId) + msg
	}
	if b.stack {
		v = append(v, TakeStacktrace(b.depth+1))
	}
	event := b.Event()
	if b.depth != 0 {
		lineDesc, fn := codeline(b.depth)
		event = event.Str("codeline", lineDesc)
		event = event.Str("func", fn)
	}
	if b.err != nil {
		event = event.Err(b.err)
	}
	event.Msgf(msg, v...)
}

func codeline(n int) (lineDesc string, fn string) {
	funcName, file, line, ok := runtime.Caller(n)
	if !ok {
		return "", ""
	}
	if i := strings.Index(file, "src/"); i >= 0 {
		return fmt.Sprintf("%s:%d", file[i+4:], line), runtime.FuncForPC(funcName).Name()
	}
	return fmt.Sprintf("%s:%d", file, line), runtime.FuncForPC(funcName).Name()
}
