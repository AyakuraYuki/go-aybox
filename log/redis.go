package log

import (
	"context"
	"fmt"
	"io"
	"time"

	"code.cloudfoundry.org/go-diodes"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

var _ zerolog.LevelWriter = (*RedisWriter)(nil)

type RedisWriter struct {
	level     zerolog.Level
	redisURL  string
	redisAuth string
	logKey    string
	client    *redis.Client
}

type RedisWriterOption func(*RedisWriter)

func NewRedisWriter(opts ...RedisWriterOption) io.Writer {
	writer := &RedisWriter{
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

	if async {
		asyncWriter := NewAsyncWriter(writer.level, writer, diodes.NewManyToOne(1024, diodes.AlertFunc(func(missed int) {
			fmt.Printf("redis writer dropped %d messages\n", missed)
		})), 1*time.Second)
		registerCloseFn(asyncWriter.Close)
		return asyncWriter
	}

	return writer
}

// Write writes data to writer
func (c *RedisWriter) Write(p []byte) (n int, err error) {
	return len(p), c.client.LPush(context.Background(), c.logKey, p).Err()
}

// WriteLevel writes data to writer with level info provided
func (c *RedisWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	if level < c.level {
		return len(p), nil
	}
	return c.Write(p)
}

func WithRedisLogLevel(level zerolog.Level) RedisWriterOption {
	return func(c *RedisWriter) {
		c.level = level
	}
}

func WithRedisURL(redisURL string) RedisWriterOption {
	return func(c *RedisWriter) {
		if redisURL != "" {
			c.redisURL = redisURL
		}
	}
}

func WithRedisAuth(redisAuth string) RedisWriterOption {
	return func(c *RedisWriter) {
		c.redisAuth = redisAuth
	}
}

func WithRedisLogKey(logKey string) RedisWriterOption {
	return func(c *RedisWriter) {
		if logKey != "" {
			c.logKey = logKey
		}
	}
}
