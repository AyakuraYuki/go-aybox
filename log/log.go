package log

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/rs/zerolog"

	"github.com/AyakuraYuki/go-aybox/ip"
)

var (
	logger = func() *zerolog.Logger {
		l := zerolog.New(NewConsoleWriter()).With().Timestamp().Str("host", ip.Hostname()).Logger()
		return &l
	}()
	callDepth   = 2
	async       = false
	defaultName = "ay-zlog"
)

func init() {
	zerolog.MessageFieldName = "desc"
	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimeFieldFormat = ""
}

type Log struct {
	depth   int
	level   zerolog.Level
	stack   bool
	err     error
	traceId string
	zlogger *zerolog.Logger
	sampler zerolog.Sampler
}

func NewLog() *Log {
	return defaultLogger(nil, callDepth, zerolog.DebugLevel, defaultName)
}

func NewLogger(level zerolog.Level, project, env string, redisOpts ...string) *Log {
	l := defaultLogger(nil, callDepth, level, project)
	WithAsync()
	WithAttachment(map[string]string{
		"project": project,
		"env":     env,
	})

	w := append([]io.Writer{}, NewConsoleWriter())
	if len(redisOpts) > 0 && redisOpts[0] != "" {
		// redis url detected
		switch len(redisOpts) {
		case 1:
			w = append(w, NewRedisWriter(
				WithRedisLogLevel(level),
				WithRedisURL(redisOpts[0])))
		case 2:
			w = append(w, NewRedisWriter(
				WithRedisLogLevel(level),
				WithRedisURL(redisOpts[0]),
				WithRedisAuth(redisOpts[1])))
		}
	}

	WithOutput(w...)
	return l
}

func defaultLogger(b *Log, depth int, level zerolog.Level, name ...string) *Log {
	if b == nil {
		b = &Log{zlogger: logger}
	} else {
		// snapshot
		b = func() *Log {
			nb := *b
			return &nb
		}()
	}
	b.depth = depth
	b.level = level
	if len(name) > 0 {
		fields := map[string]any{
			"name": name[0],
		}
		b.zlogger = logPtr(b.zlogger.With().Fields(fields).Logger())
	}
	return b
}

func Debug(name ...string) *Log {
	return defaultLogger(nil, callDepth, zerolog.DebugLevel, name...)
}

func Info(name ...string) *Log {
	return defaultLogger(nil, callDepth, zerolog.InfoLevel, name...)
}

func Warn(name ...string) *Log {
	return defaultLogger(nil, callDepth, zerolog.WarnLevel, name...)
}

func Error(name ...string) *Log {
	return defaultLogger(nil, callDepth, zerolog.ErrorLevel, name...)
}

// Context add wrapped *Log to context
func (b *Log) Context() context.Context {
	return context.WithValue(context.Background(), contextKey{}, &Context{logger: b})
}

// KV is log kv pairs.
func (b *Log) KV(key string, val string) *Log {
	b.zlogger = logPtr(b.zlogger.With().Str(key, val).Logger())
	return b
}

func (b *Log) TraceID(traceId string) *Log {
	b.traceId = traceId
	return b
}

// Stack enables stack trace
func (b *Log) Stack() *Log {
	b.stack = true
	return b
}

func (b *Log) Err(err error) *Log {
	if b.level == zerolog.ErrorLevel {
		b.err = err
	}
	return b
}

// Event returns zerolog.Event which contains Log details
func (b *Log) Event() *zerolog.Event {
	var l = *b.zlogger
	if b.sampler != nil {
		l = l.Sample(b.sampler)
	}
	return l.WithLevel(b.level)
}

// Msg output
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
		b.Msgf(fmtStr[:len(fmtStr)-2], msg...) // shrink last ', '
	}
}

// Msgf formatted output
func (b *Log) Msgf(msg string, v ...any) {

	if b.depth != 0 {
		msg = codeline(b.depth) + msg
	}

	if b.traceId != "" {
		msg = fmt.Sprintf("traceId:[%s] ", b.traceId) + msg
	}

	if b.stack {
		v = append(v, TakeStacktrace(b.depth+1))
	}

	event := b.Event()

	if b.err != nil {
		event = event.Err(b.err)
	}

	event.Msgf(msg, v...)
}

// 实现 zk logger

// Printf logs to INFO log. Arguments are handled in the manner of fmt.Printf.
func (b *Log) Printf(format string, args ...any) {
	b.Msgf(format, args)
}

// 实现 grpclog v2

// Info logs to INFO log. Arguments are handled in the manner of fmt.Print.
func (b *Log) Info(args ...any) {
	b.Msg(args)
}

// Infoln logs to INFO log. Arguments are handled in the manner of fmt.Println.
func (b *Log) Infoln(args ...any) {
	b.Msg(args)
}

// Infof logs to INFO log. Arguments are handled in the manner of fmt.Printf.
func (b *Log) Infof(format string, args ...any) {
	b.Msgf(format, args)
}

// Warning logs to WARNING log. Arguments are handled in the manner of fmt.Print.
func (b *Log) Warning(args ...any) {
	b.level = zerolog.WarnLevel
	b.Msg(args)
}

// Warningln logs to WARNING log. Arguments are handled in the manner of fmt.Println.
func (b *Log) Warningln(args ...any) {
	b.level = zerolog.WarnLevel
	b.Msg(args)
}

// Warningf logs to WARNING log. Arguments are handled in the manner of fmt.Printf.
func (b *Log) Warningf(format string, args ...any) {
	b.level = zerolog.WarnLevel
	b.Msgf(format, args)
}

// Error logs to ERROR log. Arguments are handled in the manner of fmt.Print.
func (b *Log) Error(args ...any) {
	b.level = zerolog.ErrorLevel
	b.Msg(args)
}

// Errorln logs to ERROR log. Arguments are handled in the manner of fmt.Println.
func (b *Log) Errorln(args ...any) {
	b.level = zerolog.ErrorLevel
	b.Msg(args)
}

// Errorf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
func (b *Log) Errorf(format string, args ...any) {
	b.level = zerolog.ErrorLevel
	b.Msgf(format, args)
}

// Fatal logs to ERROR log. Arguments are handled in the manner of fmt.Print.
// gRPC ensures that all Fatal logs will exit with os.Exit(1).
// Implementations may also call os.Exit() with a non-zero exit code.
func (b *Log) Fatal(args ...any) {
	b.level = zerolog.FatalLevel
	b.Msg(args)
}

// Fatalln logs to ERROR log. Arguments are handled in the manner of fmt.Println.
// gRPC ensures that all Fatal logs will exit with os.Exit(1).
// Implementations may also call os.Exit() with a non-zero exit code.
func (b *Log) Fatalln(args ...any) {
	b.level = zerolog.FatalLevel
	b.Msg(args)
}

// Fatalf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
// gRPC ensures that all Fatal logs will exit with os.Exit(1).
// Implementations may also call os.Exit() with a non-zero exit code.
func (b *Log) Fatalf(format string, args ...any) {
	b.level = zerolog.FatalLevel
	b.Msgf(format, args)
}

// V reports whether verbosity level l is at least the requested verbose level.
func (b *Log) V(l int) bool {
	return true
}

func codeline(n int) string {
	funcName, file, line, ok := runtime.Caller(n)
	if !ok {
		return ""
	}
	if i := strings.Index(file, "src/"); i != 0 {
		return "[" + file[i+4:] + ":" + strconv.Itoa(line) + " " + runtime.FuncForPC(funcName).Name() + "] "
	}
	return "[" + file + ":" + strconv.Itoa(line) + " " + runtime.FuncForPC(funcName).Name() + "] "
}

func logPtr(z zerolog.Logger) *zerolog.Logger { return &z }

func runInK8S() bool {
	return strings.EqualFold(os.Getenv("LogEnv"), "k8s")
}

type CloseFn func() error

var closeFuncs []CloseFn

func registerCloseFn(fn CloseFn) {
	closeFuncs = append(closeFuncs, fn)
}

// Close all log writer.
func Close() {
	for i := range closeFuncs {
		_ = closeFuncs[i]()
	}
}
