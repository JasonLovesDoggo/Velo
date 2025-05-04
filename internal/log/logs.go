package log

import (
	"context"
	"io"
	"log/slog"
	"os"
	"sync"
)

var (
	defaultLogger *slog.Logger
	once          sync.Once
)

func init() {
	defaultLogger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

// SetOutput sets the output destination for the default logger
func SetOutput(w io.Writer) {
	defaultLogger = slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

// SetLevel sets the logging level for the default logger
func SetLevel(level slog.Level) {
	defaultLogger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
}

func Debug(msg string, args ...any) {
	defaultLogger.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	defaultLogger.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	defaultLogger.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	defaultLogger.Error(msg, args...)
}

func Fatal(msg string, args ...any) {
	Error(msg, args...)
	os.Exit(1)
}
func With(args ...any) *slog.Logger {
	return defaultLogger.With(args...)
}

func WithContext(ctx context.Context) *slog.Logger {
	return defaultLogger.With("context", ctx)
}

func GetLogger() *slog.Logger {
	return defaultLogger
}
