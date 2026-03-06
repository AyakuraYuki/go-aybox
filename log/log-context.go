package log

import "github.com/rs/zerolog"

type contextKey struct{}

type Context struct {
	logger *Log
}

func Sample(sampler zerolog.Sampler) *Context {
	l := defaultLogger(nil, callDepth, zerolog.Disabled)
	l.sampler = sampler
	return &Context{logger: l}
}

// Debug level msg
func (b *Context) Debug(name ...string) *Log {
	return defaultLogger(b.logger, callDepth, zerolog.DebugLevel, name...)
}

// Info level msg
func (b *Context) Info(name ...string) *Log {
	return defaultLogger(b.logger, callDepth, zerolog.InfoLevel, name...)
}

// Warn level msg
func (b *Context) Warn(name ...string) *Log {
	return defaultLogger(b.logger, callDepth, zerolog.WarnLevel, name...)
}

// Error level msg
func (b *Context) Error(name ...string) *Log {
	return defaultLogger(b.logger, callDepth, zerolog.ErrorLevel, name...)
}

// KV is log kv pairs.
func (b *Context) KV(key string, val string) *Log {
	b.logger.KV(key, val)
	return b.logger
}
