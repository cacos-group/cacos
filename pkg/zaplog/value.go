package zaplog

import (
	"context"
	"go.uber.org/zap"
	"runtime"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/otel/trace"
)

var (
	// DefaultCaller is a Valuer that returns the file and line.
	DefaultCaller = Caller(3)
)

// Valuer is returns a zaplog value.
type Valuer func(ctx context.Context) interface{}

// Value return the function value.
func Value(ctx context.Context, v interface{}) interface{} {
	if v, ok := v.(Valuer); ok {
		return v(ctx)
	}
	return v
}

// Caller returns returns a Valuer that returns a pkg/file:line description of the caller.
func Caller(depth int) Valuer {
	return func(context.Context) interface{} {
		_, file, line, _ := runtime.Caller(depth)
		if strings.LastIndex(file, "/zaplog/filter.go") > 0 {
			depth++
			_, file, line, _ = runtime.Caller(depth)
		}
		if strings.LastIndex(file, "/zaplog/log.go") > 0 {
			depth++
			_, file, line, _ = runtime.Caller(depth)
		}
		idx := strings.LastIndexByte(file, '/')
		return file[idx+1:] + ":" + strconv.Itoa(line)
	}
}

// Timestamp returns a timestamp Valuer with a custom time format.
func Timestamp(layout string) Valuer {
	return func(context.Context) interface{} {
		return time.Now().Format(layout)
	}
}

// TraceID returns a traceid valuer.
func TraceID() Valuer {
	return func(ctx context.Context) interface{} {
		if span := trace.SpanContextFromContext(ctx); span.HasTraceID() {
			return span.TraceID().String()
		}
		return ""
	}
}

// SpanID returns a spanid valuer.
func SpanID() Valuer {
	return func(ctx context.Context) interface{} {
		if span := trace.SpanContextFromContext(ctx); span.HasSpanID() {
			return span.SpanID().String()
		}
		return ""
	}
}

func bindValues(ctx context.Context, keyvals []interface{}) (fields []zap.Field) {
	fields = make([]zap.Field, 0, len(keyvals)/2)

	for i := 1; i < len(keyvals); i += 2 {
		if v, ok := keyvals[i].(Valuer); ok {
			fields = append(fields, Any(keyvals[i-1].(string), v(ctx)))
		}
	}
	return fields
}

func containsValuer(keyvals []interface{}) bool {
	for i := 1; i < len(keyvals); i += 2 {
		if _, ok := keyvals[i].(Valuer); ok {
			return true
		}
	}
	return false
}
