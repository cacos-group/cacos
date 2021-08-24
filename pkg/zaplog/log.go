package zaplog

import (
	"context"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

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

type logger struct {
	log    *zap.Logger
	filter *Filter

	prefixMap map[string]interface{}
	prefix    []interface{}
	hasValuer bool
	ctx       context.Context
}

func newLogger(cfg Config) (l *logger) {
	var ws []zapcore.WriteSyncer
	if cfg.LogName != "" {
		ws = append(ws, zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.LogName,
			MaxSize:    cfg.MaxSize,    //MaxSize：在进行切割之前，日志文件的最大大小（以MB为单位）
			MaxBackups: cfg.MaxBackups, //MaxBackups：保留旧文件的最大个数
			MaxAge:     cfg.MaxAge,     //MaxAges：保留旧文件的最大天数
			Compress:   cfg.Compress,   //Compress：是否压缩/归档旧文件
		}))
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
				EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
				EncodeDuration: zapcore.SecondsDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			}),
			zapcore.NewMultiWriteSyncer(ws...),
			zapcore.DebugLevel))

	var (
		filterKeys   []string
		filterValues []string
	)

	filter := NewFilter(
		FilterLevel(LevelDebug),
		FilterKey(filterKeys...),
		FilterValue(filterValues...),
	)

	l = &logger{
		log:    zapLogger,
		filter: filter,
	}
	l.with(
		"app_id", "asdas",
		"trace_id", TraceID(),
		"span_id", SpanID(),
		"caller", DefaultCaller)

	return
}

func (l *logger) Debug(msg string, fields ...zap.Field) {
	ctx := l.ctx
	fields = append(bindValues(ctx, l.prefix), fields...)
	l.filter.formatFields(fields...)

	l.log.Debug(msg, fields...)
}
func (l *logger) Info(msg string, fields ...zap.Field) {
	ctx := l.ctx
	fields = append(bindValues(ctx, l.prefix), fields...)
	l.log.Info(msg, fields...)
}
func (l *logger) Warn(msg string, fields ...zap.Field) {
	ctx := l.ctx
	fields = append(bindValues(ctx, l.prefix), fields...)
	l.log.Warn(msg, fields...)
}
func (l *logger) Error(msg string, fields ...zap.Field) {
	ctx := l.ctx
	fields = append(bindValues(ctx, l.prefix), fields...)
	l.log.Error(msg, fields...)
}
func (l *logger) DPanic(msg string, fields ...zap.Field) {
	ctx := l.ctx
	fields = append(bindValues(ctx, l.prefix), fields...)
	l.log.DPanic(msg, fields...)
}
func (l *logger) Panic(msg string, fields ...zap.Field) {
	ctx := l.ctx
	fields = append(bindValues(ctx, l.prefix), fields...)
	l.log.Panic(msg, fields...)
}
func (l *logger) Fatal(msg string, fields ...zap.Field) {
	ctx := l.ctx
	fields = append(bindValues(ctx, l.prefix), fields...)
	l.log.Fatal(msg, fields...)
}

func (l *logger) Debugc(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(bindValues(ctx, l.prefix), fields...)
	l.filter.formatFields(fields...)

	l.log.Debug(msg, fields...)
}

func (l *logger) Infoc(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(bindValues(ctx, l.prefix), fields...)
	l.log.Info(msg, fields...)
}

func (l *logger) Warnc(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(bindValues(ctx, l.prefix), fields...)
	l.log.Warn(msg, fields...)
}

func (l *logger) Errorc(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(bindValues(ctx, l.prefix), fields...)
	l.log.Error(msg, fields...)
}

func (l *logger) DPanicc(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(bindValues(ctx, l.prefix), fields...)
	l.log.DPanic(msg, fields...)
}

func (l *logger) Panicc(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(bindValues(ctx, l.prefix), fields...)
	l.log.Panic(msg, fields...)
}

func (l *logger) Fatalc(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(bindValues(ctx, l.prefix), fields...)
	l.log.Fatal(msg, fields...)
}

// with with logger fields.
func (l *logger) with(kv ...interface{}) {
	kvs := make([]interface{}, 0, len(l.prefix)+len(kv))
	kvs = append(kvs, kv...)
	kvs = append(kvs, l.prefix...)

	l.prefix = kvs
	l.hasValuer = containsValuer(kvs)
}
