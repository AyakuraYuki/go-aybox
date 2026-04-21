package gorm_log

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go/ext"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"

	"github.com/AyakuraYuki/go-aybox/log"
	ayopentracing "github.com/AyakuraYuki/go-aybox/opentracing"
)

const (
	dbLogName = "DBLog"
	dbSQLSlow = "DBSQLSlow"
)

type LogLevel int

const (
	LogLevelSilent LogLevel = iota + 1
	LogLevelError
	LogLevelWarn
	LogLevelInfo
)

var _ logger.Interface = (*dbLog)(nil)

type dbLog struct {
	logger.Config
	zl           *log.Logger
	infoStr      string
	warnStr      string
	errStr       string
	traceStr     string
	traceErrStr  string
	traceWarnStr string
}

// New creates a gorm logger.Interface backed by the provided *log.Logger.
func New(config logger.Config, zl *log.Logger) logger.Interface {
	var (
		infoStr      = "%s "
		warnStr      = "%s "
		errStr       = "%s "
		traceStr     = "%s [%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s [%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s [%.3fms] [rows:%v] %s"
	)

	return &dbLog{
		Config:       config,
		zl:           zl,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
}

// LogMode log mode
func (l *dbLog) LogMode(level logger.LogLevel) logger.Interface {
	nl := *l
	nl.LogLevel = level
	return &nl
}

// Info print info
func (l *dbLog) Info(ctx context.Context, msg string, data ...any) {
	if l.LogLevel >= logger.Info {
		l.zl.InfoL(dbLogName).Msgf(l.infoStr+msg, append([]any{utils.FileWithLineNum()}, data...)...)
	}
}

// Warn print warn messages
func (l *dbLog) Warn(ctx context.Context, msg string, data ...any) {
	if l.LogLevel >= logger.Warn {
		l.zl.WarnL(dbLogName).Msgf(l.warnStr+msg, append([]any{utils.FileWithLineNum()}, data...)...)
	}
}

// Error print error messages
func (l *dbLog) Error(ctx context.Context, msg string, data ...any) {
	if l.LogLevel >= logger.Error {
		l.zl.ErrorL(dbLogName).Msgf(l.errStr+msg, append([]any{utils.FileWithLineNum()}, data...)...)
	}
}

// Trace print sql message
func (l *dbLog) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	codeLine := utils.FileWithLineNum()
	if arr := strings.Split(codeLine, "/cmd/"); len(arr) == 2 {
		codeLine = "/cmd/" + arr[1]
	}

	span, _ := ayopentracing.StartSpanFromContextWithSt(ctx, codeLine, begin)
	defer span.Finish()
	span.SetTag(string(ext.DBStatement), "GORM")
	span.SetTag(string(ext.DBType), "sql")
	span.LogKV("sql", sql)

	switch {
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, logger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):

		var logStr string
		if rows == -1 {
			logStr = fmt.Sprintf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			logStr = fmt.Sprintf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}

		l.zl.ErrorL(dbLogName).Msg(logStr)

	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)

		var logStr string
		if rows == -1 {
			logStr = fmt.Sprintf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			logStr = fmt.Sprintf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
		l.zl.WarnL(dbSQLSlow).Msg(logStr)

	case l.LogLevel == logger.Info:
		if rows == -1 {
			l.zl.InfoL(dbLogName).Msgf(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.zl.InfoL(dbLogName).Msgf(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}
