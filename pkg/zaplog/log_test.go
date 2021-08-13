package zaplog

import (
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
