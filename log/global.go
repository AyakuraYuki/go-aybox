package log

import (
	"context"

	"github.com/rs/zerolog"
)

// defaultLogger is the package-level singleton, ready to use out of the box.
var defaultLogger = New()

// Default returns the global Logger instance.
func Default() *Logger {
	return defaultLogger
}

// SetDefault replaces the global Logger with l.
// The old logger is NOT closed automatically; call Close on it if needed.
//
// WARNING: SetDefault is not goroutine-safe. It must only be called during
// program initialisation, before any logging goroutine is started. Calling it
// concurrently with package-level logging functions (Info, Debug, …) causes a
// data race on defaultLogger.
func SetDefault(l *Logger) {
	defaultLogger = l
}

// Configure rebuilds the global Logger with the given options and replaces it.
// The old logger is NOT closed automatically.
func Configure(opts ...Option) {
	defaultLogger = New(opts...)
}

// Debug returns a *Log at debug level from the global Logger.
func Debug(name ...string) *Log {
	return defaultLogger.DebugL(name...)
}

// Info returns a *Log at info level from the global Logger.
func Info(name ...string) *Log {
	return defaultLogger.InfoL(name...)
}

// Warn returns a *Log at warn level from the global Logger.
func Warn(name ...string) *Log {
	return defaultLogger.WarnL(name...)
}

// Error returns a *Log at error level from the global Logger.
func Error(name ...string) *Log {
	return defaultLogger.ErrorL(name...)
}

// Fatal returns a *Log at fatal level from the global Logger.
func Fatal(name ...string) *Log {
	return defaultLogger.FatalL(name...)
}

// Sample returns a *Context with sampling enabled from the global Logger.
func Sample(sampler zerolog.Sampler) *Context {
	return defaultLogger.Sample(sampler)
}

// FromContext extracts a *Context from ctx using the global Logger.
func FromContext(ctx context.Context) *Context {
	return defaultLogger.FromContext(ctx)
}

// Close flushes and closes all writers of the global Logger.
func Close() {
	defaultLogger.Close()
}
