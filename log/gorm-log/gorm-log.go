package gorm_log

import (
	"context"
	"errors"
	"fmt"
	"os"
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

var env string
var serviceName string

func init() {
	env = os.Getenv("ENV")
	if env == "" {
		env = "pro"
	}
}

func WithGormLogEnv(mEnv string)                 { env = mEnv }
func WithGormLogServiceName(mServiceName string) { serviceName = mServiceName }

type LogLevel int

const (
	LogLevelSilent LogLevel = iota + 1
	LogLevelError
	LogLevelWarn
	LogLevelInfo
)

type dbLog struct {
	logger.Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

func NewDBLog(config logger.Config) logger.Interface {
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
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info print info
func (l dbLog) Info(ctx context.Context, msg string, data ...any) {
	if l.LogLevel >= logger.Info {
		log.Info(dbLogName).Msgf(l.infoStr+msg, append([]any{utils.FileWithLineNum()}, data...)...)
	}
}

// Warn print warn messages
func (l dbLog) Warn(ctx context.Context, msg string, data ...any) {
	if l.LogLevel >= logger.Warn {
		log.Warn(dbLogName).Msgf(l.warnStr+msg, append([]any{utils.FileWithLineNum()}, data...)...)
	}
}

// Error print error messages
func (l dbLog) Error(ctx context.Context, msg string, data ...any) {
	if l.LogLevel >= logger.Error {
		log.Error(dbLogName).Msgf(l.errStr+msg, append([]any{utils.FileWithLineNum()}, data...)...)
	}
}

// Trace print sql message
func (l dbLog) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
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

		log.Error(dbLogName).Msg(logStr)

	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)

		var logStr string
		if rows == -1 {
			logStr = fmt.Sprintf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			logStr = fmt.Sprintf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
		log.Warn(dbSQLSlow).Msgf(logStr)

	case l.LogLevel == logger.Info:
		if rows == -1 {
			log.Info(dbLogName).Msgf(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			log.Info(dbLogName).Msgf(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}
