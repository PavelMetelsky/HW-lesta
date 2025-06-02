package logger

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

// Выполнено на основе https://github.com/orandin/slog-gorm/

var GormLogger *DbGormLogger

type LogType string

const (
	ErrorLogType     LogType = "sql_error"
	SlowQueryLogType LogType = "slow_query"
	DefaultLogType   LogType = "default"

	SourceField    = "file"
	ErrorField     = "error"
	QueryField     = "query"
	DurationField  = "duration"
	SlowQueryField = "slow_query"
	RowsField      = "rows"
)

type DbGormLogger struct {
	slogger                   *slog.Logger
	ignoreTrace               bool
	ignoreRecordNotFoundError bool
	traceAll                  bool
	slowThreshold             time.Duration
	logLevel                  map[LogType]slog.Level
	sourceField               string
	errorField                string
	Context                   *context.Context
}

func (logger DbGormLogger) LogMode(_ gormlogger.LogLevel) gormlogger.Interface {
	return logger
}

func (logger DbGormLogger) Info(ctx context.Context, msg string, args ...any) {
	logger.slogger.InfoContext(ctx, msg, args...)
}

func (logger DbGormLogger) Warn(ctx context.Context, msg string, args ...any) {
	logger.slogger.WarnContext(ctx, msg, args...)
}

func (logger DbGormLogger) Error(ctx context.Context, msg string, args ...any) {
	logger.slogger.ErrorContext(ctx, msg, args...)
}

// Trace Регистрация событий (выполнение запросов к базе)
func (logger DbGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if logger.ignoreTrace {
		return
	}
	elapsed := time.Since(begin)

	switch {
	case err != nil && (!errors.Is(err, gorm.ErrRecordNotFound) || !logger.ignoreRecordNotFoundError):
		sql, rows := fc()

		logger.slogger.Log(ctx, logger.logLevel[ErrorLogType], err.Error(),
			slog.Any(logger.errorField, err),
			slog.String(QueryField, sql),
			slog.Duration(DurationField, elapsed),
			slog.Int64(RowsField, rows),
			slog.String(logger.sourceField, utils.FileWithLineNum()),
		)
	case logger.slowThreshold != 0 && elapsed > logger.slowThreshold:
		sql, rows := fc()

		logger.slogger.Log(ctx, logger.logLevel[SlowQueryLogType], fmt.Sprintf("slow sql query [%s >= %v]", elapsed, logger.slowThreshold),
			slog.Bool(SlowQueryField, true),
			slog.String(QueryField, sql),
			slog.Duration(DurationField, elapsed),
			slog.Int64(RowsField, rows),
			slog.String(logger.sourceField, utils.FileWithLineNum()),
		)

	case logger.traceAll:
		sql, rows := fc()
		logger.slogger.Log(ctx, logger.logLevel[DefaultLogType], fmt.Sprintf("SQL query executed [%s]", elapsed),
			slog.String(QueryField, sql),
			slog.Duration(DurationField, elapsed),
			slog.Int64(RowsField, rows),
			slog.String(logger.sourceField, utils.FileWithLineNum()),
		)
	}
}

func CreateDbLogger() *DbGormLogger {
	optsApi := slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	slowThreshold := 30 * time.Millisecond
	logger := DbGormLogger{
		slogger:                   slog.New(slog.NewJSONHandler(os.Stdout, &optsApi)),
		ignoreRecordNotFoundError: true,
		errorField:                ErrorField,
		sourceField:               SourceField,
		traceAll:                  true,
		slowThreshold:             slowThreshold,

		logLevel: map[LogType]slog.Level{
			ErrorLogType:     slog.LevelError,
			SlowQueryLogType: slog.LevelWarn,
			DefaultLogType:   slog.LevelInfo,
		},
	}
	GormLogger = &logger
	return GormLogger
}
