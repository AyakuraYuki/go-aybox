package log

import (
	"context"
	"fmt"
	"io"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

var _ hlog.FullLogger = (*Logger)(nil)

func (l *Logger) SetLevel(_ hlog.Level) {}
func (l *Logger) SetOutput(_ io.Writer) {}

func (l *Logger) Trace(v ...any)  { l.DebugL("hertz-Trace").Msg(fmt.Sprint(v...)) }
func (l *Logger) Debug(v ...any)  { l.DebugL("hertz-Debug").Msg(fmt.Sprint(v...)) }
func (l *Logger) Info(v ...any)   { l.InfoL("hertz-Info").Msg(fmt.Sprint(v...)) }
func (l *Logger) Notice(v ...any) { l.InfoL("hertz-Notice").Msg(fmt.Sprint(v...)) }
func (l *Logger) Warn(v ...any)   { l.WarnL("hertz-Warn").Msg(fmt.Sprint(v...)) }
func (l *Logger) Error(v ...any)  { l.ErrorL("hertz-Error").Msg(fmt.Sprint(v...)) }
func (l *Logger) Fatal(v ...any)  { l.FatalL("hertz-Fatal").Msg(fmt.Sprint(v...)) }

func (l *Logger) Tracef(format string, v ...any)  { l.DebugL("hertz-Trace").Msgf(format, v...) }
func (l *Logger) Debugf(format string, v ...any)  { l.DebugL("hertz-Debug").Msgf(format, v...) }
func (l *Logger) Infof(format string, v ...any)   { l.InfoL("hertz-Info").Msgf(format, v...) }
func (l *Logger) Noticef(format string, v ...any) { l.InfoL("hertz-Notice").Msgf(format, v...) }
func (l *Logger) Warnf(format string, v ...any)   { l.WarnL("hertz-Warn").Msgf(format, v...) }
func (l *Logger) Errorf(format string, v ...any)  { l.ErrorL("hertz-Error").Msgf(format, v...) }
func (l *Logger) Fatalf(format string, v ...any)  { l.FatalL("hertz-Fatal").Msgf(format, v...) }

func (l *Logger) CtxTracef(_ context.Context, format string, v ...any) {
	l.DebugL("hertz-Trace").Msgf(format, v...)
}

func (l *Logger) CtxDebugf(_ context.Context, format string, v ...any) {
	l.DebugL("hertz-Debug").Msgf(format, v...)
}

func (l *Logger) CtxInfof(_ context.Context, format string, v ...any) {
	l.InfoL("hertz-Info").Msgf(format, v...)
}

func (l *Logger) CtxNoticef(_ context.Context, format string, v ...any) {
	l.InfoL("hertz-Notice").Msgf(format, v...)
}

func (l *Logger) CtxWarnf(_ context.Context, format string, v ...any) {
	l.WarnL("hertz-Warn").Msgf(format, v...)
}

func (l *Logger) CtxErrorf(_ context.Context, format string, v ...any) {
	l.ErrorL("hertz-Error").Msgf(format, v...)
}

func (l *Logger) CtxFatalf(_ context.Context, format string, v ...any) {
	l.FatalL("hertz-Fatal").Msgf(format, v...)
}
