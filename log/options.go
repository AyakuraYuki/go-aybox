package log

import (
	"io"

	"github.com/rs/zerolog"
)

// config holds the builder state used by New.
type config struct {
	level    zerolog.Level
	depth    int
	async    bool
	codeline bool
	writers  []io.Writer
	fields   map[string]string
	hostname string
}

// Option configures a Logger at construction time.
type Option func(*config)

// WithLevel sets the minimum log level and the zerolog global level.
func WithLevel(level zerolog.Level) Option {
	return func(c *config) {
		c.level = level
		zerolog.SetGlobalLevel(level)
	}
}

// WithDepth overrides the default call-stack depth used for source-line
// annotation.
func WithDepth(depth int) Option {
	return func(c *config) {
		c.depth = depth
	}
}

// WithAsync enables non-blocking, lock-free log writing via a diode buffer.
// Call Logger.Close to flush pending entries before process exit.
func WithAsync() Option {
	return func(c *config) {
		c.async = true
	}
}

// WithCodeline enables file:line ("codeline") and function name ("func") fields
// on every log entry. Disabled by default.
func WithCodeline() Option {
	return func(c *config) {
		c.codeline = true
	}
}

// WithWriters adds one or more writers.
// Multiple calls are additive; the first call overrides the default
// ConsoleWriter.
func WithWriters(w ...io.Writer) Option {
	return func(c *config) {
		c.writers = append(c.writers, w...)
	}
}

// WithFields adds static key-value string fields to every log entry.
func WithFields(kv map[string]string) Option {
	return func(c *config) {
		if c.fields == nil {
			c.fields = make(map[string]string)
		}
		for k, v := range kv {
			c.fields[k] = v
		}
	}
}

// WithHostname changes hostname.
func WithHostname(hostname string) Option {
	return func(c *config) {
		c.hostname = hostname
	}
}
