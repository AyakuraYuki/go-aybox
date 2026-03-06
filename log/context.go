package log

import (
	"context"

	"github.com/rs/zerolog"
)

// Sample returns a *Context with the given sampler applied to every entry.
func (l *Logger) Sample(sampler zerolog.Sampler) *Context {
	return &Context{logger: l, sampler: sampler}
}

// FromContext extracts a *Context stored by (*Context).ToContext.
// Falls back to a plain *Context backed by this logger if none is found.
func (l *Logger) FromContext(ctx context.Context) *Context {
	if lc, ok := ctx.Value(contextKey{}).(*Context); ok {
		return lc
	}
	return &Context{logger: l}
}

// Context carries shared log state (sampler, extra fields) across multiple log
// entries.
// Obtain one via Logger.Sample or Logger.FromContext.
type Context struct {
	logger  *Logger
	zlogger *zerolog.Logger // optional override; nil means use logger.zlogger
	sampler zerolog.Sampler
}

type contextKey struct{}

// Debug returns a *Log at debug level.
// Returns nil when c is nil, mirroring the nil *Log fast-path.
func (c *Context) Debug(name ...string) *Log {
	if c == nil {
		return nil
	}
	return c.newLog(zerolog.DebugLevel, name...)
}

// Info returns a *Log at info level.
// Returns nil when c is nil, mirroring the nil *Log fast-path.
func (c *Context) Info(name ...string) *Log {
	if c == nil {
		return nil
	}
	return c.newLog(zerolog.InfoLevel, name...)
}

// Warn returns a *Log at warn level.
// Returns nil when c is nil, mirroring the nil *Log fast-path.
func (c *Context) Warn(name ...string) *Log {
	if c == nil {
		return nil
	}
	return c.newLog(zerolog.WarnLevel, name...)
}

// Error returns a *Log at error level.
// Returns nil when c is nil, mirroring the nil *Log fast-path.
func (c *Context) Error(name ...string) *Log {
	if c == nil {
		return nil
	}
	return c.newLog(zerolog.ErrorLevel, name...)
}

// KV returns a new Context with an additional string field shared across all
// log entries created from it. Returns nil when c is nil.
// This is setup-time and not in the hot path.
func (c *Context) KV(key, val string) *Context {
	if c == nil {
		return nil
	}
	base := c.baseLogger()
	zl := base.With().Str(key, val).Logger()
	nc := *c
	nc.zlogger = &zl
	return &nc
}

// ToContext stores this Context in the provided Go context.
// Retrieve it later with (*Logger).FromContext.
// Returns parent unchanged when c is nil.
func (c *Context) ToContext(parent context.Context) context.Context {
	if c == nil {
		return parent
	}
	return context.WithValue(parent, contextKey{}, c)
}

func (c *Context) newLog(level zerolog.Level, name ...string) *Log {
	if c == nil {
		return nil
	}
	zl := *c.baseLogger()
	if c.sampler != nil {
		zl = zl.Sample(c.sampler)
	}
	event := zl.WithLevel(level)
	if event == nil {
		return nil
	}
	if len(name) > 0 && name[0] != "" {
		event = event.Str("name", name[0])
	}
	b := logPool.Get().(*Log)
	b.event = event
	b.depth = c.logger.depth
	b.level = level
	b.stack = false
	b.traceId = ""
	b.logger = c.logger
	return b
}

func (c *Context) baseLogger() *zerolog.Logger {
	if c.zlogger != nil {
		return c.zlogger
	}
	return &c.logger.zlogger
}
