package log

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"code.cloudfoundry.org/go-diodes"
	"github.com/rs/zerolog"

	"github.com/AyakuraYuki/go-aybox/log/async"
	"github.com/AyakuraYuki/go-aybox/log/console"
)

// Leveler is implemented by writers that expose a minimum log level.
type Leveler interface {
	Level() zerolog.Level
}

// Logger is a multi-instance logger backed by zerolog.
// Logging methods (Debug, Info, Warn, Error) are safe to call concurrently.
// With and Accept are intended for setup time; Accept is mutex-protected to
// serialize concurrent writes to the underlying logger.
// Create one with New().
type Logger struct {
	zlogger  zerolog.Logger
	mu       sync.RWMutex // guards writes to zlogger via Accept
	depth    int          // base call-stack depth for codeline; default 2
	codeline bool         // emit file:line and func fields; default false
	closeFns []CloseFn
}

// New creates a Logger with the provided options.
// If no WithWriters option is given, a ConsoleWriter is used.
func New(opts ...Option) *Logger {
	cfg := &config{
		level:    zerolog.DebugLevel,
		depth:    2,
		hostname: hostname(),
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
				time.Second,
				cfg.asyncCloseTimeout)
			wrapped[i] = aw
			// async.Writer.Close flushes the diode and then closes the
			// underlying writer if it implements io.Closer, so we only need
			// one entry per writer here.
			closeFns = append(closeFns, aw.Close)
		}
		writers = wrapped
	} else {
		// In synchronous mode, collect Close functions for any writer that
		// implements io.Closer (e.g. redis.Writer) so that Logger.Close
		// properly releases their resources.
		for _, w := range writers {
			if c, ok := w.(io.Closer); ok {
				closeFns = append(closeFns, c.Close)
			}
		}
	}

	var writer io.Writer
	if len(writers) == 1 {
		writer = writers[0]
	} else {
		writer = zerolog.MultiLevelWriter(writers...)
	}

	logger := zerolog.New(writer).
		Level(cfg.level).
		With().
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

func hostname() string {
	name, err := os.Hostname()
	if err == nil {
		return name
	}

	f, err := os.Open("/etc/hostname")
	if err != nil {
		return ""
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	r, err := io.ReadAll(f)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(r))
}

// With exposes the underlying zerolog.Context for direct field configuration.
// Intended for setup time before the logger is shared across goroutines.
// Commit changes back with Accept.
func (l *Logger) With() zerolog.Context {
	return l.zlogger.With()
}

// Accept applies a configured zerolog.Context to this logger.
// A write lock is held for the duration to serialise concurrent Accept calls.
// Intended for setup time; do not interleave Accept with active logging.
func (l *Logger) Accept(ctx zerolog.Context) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.zlogger = ctx.Logger()
}

// Log starts a new message with no level.
func (l *Logger) Log(name ...string) *Log { return l.newLog(zerolog.NoLevel, name...) }

// DebugL returns a *Log at debug level, optionally tagged with name.
func (l *Logger) DebugL(name ...string) *Log { return l.newLog(zerolog.DebugLevel, name...) }

// InfoL returns a *Log at info level, optionally tagged with name.
func (l *Logger) InfoL(name ...string) *Log { return l.newLog(zerolog.InfoLevel, name...) }

// WarnL returns a *Log at warn level, optionally tagged with name.
func (l *Logger) WarnL(name ...string) *Log { return l.newLog(zerolog.WarnLevel, name...) }

// ErrorL returns a *Log at error level, optionally tagged with name.
func (l *Logger) ErrorL(name ...string) *Log { return l.newLog(zerolog.ErrorLevel, name...) }

// FatalL returns a *Log at fatal level, optionally tagged with name.
func (l *Logger) FatalL(name ...string) *Log { return l.newLog(zerolog.FatalLevel, name...) }

// logPool recycles *Log objects to reduce per-call heap allocations.
// Objects are only borrowed from the pool when a level is enabled; the nil
// fast-path for disabled levels never touches the pool.
//
// Contract: after Msg or Msgf returns, the caller must not retain or use the
// *Log — it is returned to the pool immediately after emission.
var logPool = sync.Pool{
	New: func() any {
		return new(Log)
	},
}

// newLog creates a *Log for the given level.
// Returns nil without allocating when the level is disabled — same fast path
// that zerolog uses for a nil *Event. Otherwise, borrows a *Log from logPool
// and initializes every field to avoid stale state from a prior use.
func (l *Logger) newLog(level zerolog.Level, name ...string) *Log {
	event := l.zlogger.WithLevel(level)
	if event == nil {
		return nil
	}
	if len(name) > 0 && name[0] != "" {
		event = event.Str("name", name[0])
	}
	b := logPool.Get().(*Log)
	b.event = event
	b.depth = l.depth
	b.level = level
	b.stack = false
	b.traceId = ""
	b.logger = l
	return b
}

// CloseFn is a function that closes a resource.
type CloseFn func() error

// Close flushes and closes all writers registered to this logger.
func (l *Logger) Close() {
	for _, fn := range l.closeFns {
		_ = fn()
	}
}
