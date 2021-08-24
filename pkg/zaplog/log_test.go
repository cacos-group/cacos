package zaplog

import (
	"context"
	"testing"
)

func TestInfo(t *testing.T) {
	logger := DefaultLogger
	logger.Debug("hello world", Any("key", 1))
}

func BenchmarkLogger_Debug(b *testing.B) {
	logger := New(Config{
		LogName: "./logs/bench.log",
	})
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Debug("hello world", Any("key", 1))
	}
}

func TestLogger_Infoc(t *testing.T) {
	ctx := context.WithValue(context.Background(), "test", "1")
	DefaultLogger.Infoc(ctx, "test")
}
