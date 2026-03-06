package log

import (
	"errors"
	"io"
	"testing"
	"time"

	"github.com/rs/zerolog"
)

var (
	errExample  = errors.New("fail")
	fakeMessage = "Test logging, but use a somewhat realistic message length."
)

func BenchmarkLogEmpty(b *testing.B) {
	logger := New(WithWriters(io.Discard))
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Log().Msg("")
		}
	})
}

func BenchmarkDisabled(b *testing.B) {
	logger := New(
		WithWriters(io.Discard),
		WithLevel(zerolog.Disabled))
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info().Msg(fakeMessage)
		}
	})
}

func BenchmarkInfo(b *testing.B) {
	logger := New(WithWriters(io.Discard))
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info().Msg(fakeMessage)
		}
	})
}

func BenchmarkContextFields(b *testing.B) {
	logger := New(WithWriters(io.Discard))
	ctx := logger.With().
		Str("string", "four!").
		Time("time", time.Time{}).
		Int("int", 123).
		Float32("float", -2.203230293249593)
	logger.Accept(ctx)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info().Msg(fakeMessage)
		}
	})
}

func BenchmarkContextAppend(b *testing.B) {
	logger := New(WithWriters(io.Discard))
	ctx := logger.With().Str("foo", "bar")
	logger.Accept(ctx)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.With().Str("bar", "baz")
		}
	})
}

func BenchmarkLogFields(b *testing.B) {
	logger := New(WithWriters(io.Discard))
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info().
				Str("string", "four!").
				Time("time", time.Time{}).
				Int("int", 123).
				Float32("float", -2.203230293249593).
				Msg(fakeMessage)
		}
	})
}

type obj struct {
	Pub  string
	Tag  string `json:"tag"`
	priv int
}

func (o obj) MarshalZerologObject(e *zerolog.Event) {
	e.Str("Pub", o.Pub).
		Str("Tag", o.Tag).
		Int("priv", o.priv)
}

func BenchmarkLogArrayObject(b *testing.B) {
	obj1 := obj{"a", "b", 2}
	obj2 := obj{"c", "d", 3}
	obj3 := obj{"e", "f", 4}
	logger := New(WithWriters(io.Discard))
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		arr := zerolog.Arr()
		arr.Object(&obj1)
		arr.Object(&obj2)
		arr.Object(&obj3)
		logger.Info().Array("objects", arr).Msg("test")
	}
}
