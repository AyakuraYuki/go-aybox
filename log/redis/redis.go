package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

var _ zerolog.LevelWriter = (*Writer)(nil)
var _ interface{ Close() error } = (*Writer)(nil)

// writeTimeout is the per-call deadline applied to every LPUSH. Keeping it
// short ensures that a slow or unavailable Redis instance does not block the
// logging hot path for longer than necessary.
const writeTimeout = 5 * time.Second

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

// Write pushes p as a new entry to the Redis list identified by logKey.
// A short deadline (writeTimeout) is applied so that a slow or unavailable
// Redis instance does not block the caller indefinitely.
func (c *Writer) Write(p []byte) (n int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	defer cancel()
	// Copy p before the call: the caller may reuse the buffer immediately
	// after Write returns, and go-redis may not copy it before the network write.
	entry := make([]byte, len(p))
	copy(entry, p)
	return len(p), c.client.LPush(ctx, c.logKey, entry).Err()
}

// Close closes the underlying Redis client and releases its connection pool.
// It must be called when the writer is no longer needed; Logger.Close will
// call it automatically when this Writer is registered via WithWriters.
func (c *Writer) Close() error {
	return c.client.Close()
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
