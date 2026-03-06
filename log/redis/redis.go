package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

var _ zerolog.LevelWriter = (*Writer)(nil)

type Writer struct {
	level     zerolog.Level
	redisURL  string
	redisAuth string
	logKey    string
	client    *redis.Client
}

type Option func(*Writer)

func New(opts ...Option) *Writer {
	writer := &Writer{
		level:    zerolog.InfoLevel,
		redisURL: "redis:6379",
		logKey:   "ay:zlog:redis.writer:log",
	}
	for _, opt := range opts {
		opt(writer)
	}
	redisOpt := &redis.Options{
		Network: "tcp",
		Addr:    writer.redisURL,
	}
	if writer.redisAuth != "" {
		redisOpt.Password = writer.redisAuth
	}
	writer.client = redis.NewClient(redisOpt)
	return writer
}

// Level returns the minimum level accepted by this writer.
func (c *Writer) Level() zerolog.Level {
	return c.level
}

// Write writes data to writer
func (c *Writer) Write(p []byte) (n int, err error) {
	return len(p), c.client.LPush(context.Background(), c.logKey, p).Err()
}

// WriteLevel writes data to writer with level info provided
func (c *Writer) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	if level < c.level {
		return len(p), nil
	}
	return c.Write(p)
}

func WithLogLevel(level zerolog.Level) Option {
	return func(c *Writer) {
		c.level = level
	}
}

func WithURL(redisURL string) Option {
	return func(c *Writer) {
		if redisURL != "" {
			c.redisURL = redisURL
		}
	}
}

func WithAuth(redisAuth string) Option {
	return func(c *Writer) {
		c.redisAuth = redisAuth
	}
}

func WithLogKey(logKey string) Option {
	return func(c *Writer) {
		if logKey != "" {
			c.logKey = logKey
		}
	}
}
