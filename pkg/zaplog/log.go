package zaplog

import (
	"context"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
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
}

type Config struct {
	Level        string
	LogName      string
	Stdout       bool
	MaxSize      int  //MaxSize：在进行切割之前，日志文件的最大大小（以MB为单位）
	MaxBackups   int  //MaxBackups：保留旧文件的最大个数
	MaxAge       int  //MaxAges：保留旧文件的最大天数
	Compress     bool //Compress：是否压缩/归档旧文件
	FilterKeys   []string
	FilterValues []string
}

func New(cfg Config) Logger {
	return newLogger(cfg)
}

func newLogger(cfg Config) (l *logger) {
	l = new(logger)
	var ws []zapcore.WriteSyncer
	if cfg.LogName != "" {
		ws = append(ws, zapcore.AddSync(zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.LogName,
			MaxSize:    cfg.MaxSize,    //MaxSize：在进行切割之前，日志文件的最大大小（以MB为单位）
			MaxBackups: cfg.MaxBackups, //MaxBackups：保留旧文件的最大个数
			MaxAge:     cfg.MaxAge,     //MaxAges：保留旧文件的最大天数
			Compress:   cfg.Compress,   //Compress：是否压缩/归档旧文件
		})))
	}

	if cfg.Stdout || len(ws) == 0 {
		ws = append(ws, zapcore.AddSync(os.Stdout))
	}

	zapLogger := zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zapcore.EncoderConfig{
				TimeKey:        "ts",
				LevelKey:       "level",
				NameKey:        "logger",
				CallerKey:      "caller",
				FunctionKey:    zapcore.OmitKey,
				MessageKey:     "msg",
				StacktraceKey:  "stacktrace",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.CapitalLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.SecondsDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			}),
			zapcore.NewMultiWriteSyncer(ws...),
			zapcore.DebugLevel))

	var (
		filterKeys   []string
		filterValues []string
	)

	filterLog := NewFilter(zapLogger,
		FilterLevel(LevelDebug),
		FilterKey(filterKeys...),
		FilterValue(filterValues...),
	)

	l.log = filterLog
	l.with(
		"trace_id", TraceID(),
		"span_id", SpanID(),
		"caller", DefaultCaller)
	return
}

type logger struct {
	log       Logger
	prefixMap map[string]interface{}
	prefix    []interface{}
	hasValuer bool
	ctx       context.Context
}

func (l *logger) Debug(msg string, fields ...zap.Field) {
	l.log.Debug(msg, l.bindValues(l.ctx, fields...)...)
}
func (l *logger) Info(msg string, fields ...zap.Field) {
	l.log.Info(msg, l.bindValues(l.ctx, fields...)...)
}
func (l *logger) Warn(msg string, fields ...zap.Field) {
	l.log.Warn(msg, l.bindValues(l.ctx, fields...)...)
}
func (l *logger) Error(msg string, fields ...zap.Field) {
	l.log.Error(msg, l.bindValues(l.ctx, fields...)...)
}
func (l *logger) DPanic(msg string, fields ...zap.Field) {
	l.log.DPanic(msg, l.bindValues(l.ctx, fields...)...)
}
func (l *logger) Panic(msg string, fields ...zap.Field) {
	l.log.Panic(msg, l.bindValues(l.ctx, fields...)...)
}
func (l *logger) Fatal(msg string, fields ...zap.Field) {
	l.log.Fatal(msg, l.bindValues(l.ctx, fields...)...)
}

func (l *logger) bindValues(ctx context.Context, fields ...zap.Field) []zap.Field {
	keyvals := l.prefix
	for i := 1; i < len(keyvals); i += 2 {
		if v, ok := keyvals[i].(Valuer); ok {
			fields = append(fields, Any(keyvals[i-1].(string), v(ctx)))
		}
	}

	return fields
}

// with with logger fields.
func (l *logger) with(kv ...interface{}) {
	kvs := make([]interface{}, 0, len(l.prefix)+len(kv))
	kvs = append(kvs, kv...)
	kvs = append(kvs, l.prefix...)

	l.prefix = kvs
	l.hasValuer = containsValuer(kvs)
}
