package logger

import (
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

var ApiLog ApiGinLogger

type ApiGinLogger struct {
	Log     *slog.Logger
	Config  *gin.HandlerFunc
	Filters *gin.HandlerFunc
}

func InitApiLogger() {
	optsApi := slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &optsApi))

	configParams := sloggin.Config{
		WithRequestBody:    false,
		WithResponseBody:   false,
		WithRequestHeader:  false,
		WithResponseHeader: false,
	}
	config := sloggin.NewWithConfig(logger, configParams)

	filters := sloggin.NewWithFilters(
		logger,
		sloggin.IgnorePathPrefix("/swagger/"),
	)
	ApiLog.Log = logger
	ApiLog.Config = &config
	ApiLog.Filters = &filters
}
