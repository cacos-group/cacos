package zaplog

import (
	"context"
	"go.uber.org/zap"
)

var (
	// DefaultLogger is default logger.
	DefaultLogger Logger = New(Config{})
)

func Any(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}

// Logger is a logger interface.
type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	DPanic(msg string, fields ...zap.Field)
	Panic(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)

	Debugc(ctx context.Context, msg string, fields ...zap.Field)
	Infoc(ctx context.Context, msg string, fields ...zap.Field)
	Warnc(ctx context.Context, msg string, fields ...zap.Field)
	Errorc(ctx context.Context, msg string, fields ...zap.Field)
	DPanicc(ctx context.Context, msg string, fields ...zap.Field)
	Panicc(ctx context.Context, msg string, fields ...zap.Field)
	Fatalc(ctx context.Context, msg string, fields ...zap.Field)
}

func New(cfg Config) Logger {
	return newLogger(cfg)
}
