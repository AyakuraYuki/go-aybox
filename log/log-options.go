package log

import (
	"context"
	"io"

	"github.com/rs/zerolog"
)

// WithOutput set multiple log writers.
// Be careful, all WithXXX method are non-thread safe.
func WithOutput(w ...io.Writer) {
	switch len(w) {
	case 0:
		return
	case 1:
		*logger = logger.Output(w[0])
	default:
		*logger = logger.Output(zerolog.MultiLevelWriter(w...))
	}
}

// WithLevel set global zerolog level.
// Be careful, all WithXXX method are non-thread safe.
func WithLevel(level zerolog.Level) {
	zerolog.SetGlobalLevel(level)
}

// WithCallDepth set call depth for showing line number.
// Be careful, all WithXXX method are non-thread safe.
func WithCallDepth(depth int) {
	callDepth = depth
}

// WithAsync enables async log, should use Close function to wait all flushed.
// Be careful, all WithXXX method are non-thread safe.
func WithAsync() {
	async = true
}

// WithAttachment adds global key-value to logger.
// Be careful, all WithXXX method are non-thread safe.
func WithAttachment(kv map[string]string) {
	for k, v := range kv {
		*logger = logger.With().Str(k, v).Logger()
	}
}

// WithContext creates a wrapped Log.
// Be careful, all WithXXX method are non-thread safe.
func WithContext(ctx ...context.Context) *Context {
	if len(ctx) > 0 && ctx[0] != nil {
		if l, ok := ctx[0].Value(contextKey{}).(*Context); ok {
			return l
		}
	}
	l := defaultLogger(nil, callDepth, zerolog.Disabled)
	return &Context{logger: l}
}
