package zaplog

import "go.uber.org/zap"

// FilterOption is filter option.
type FilterOption func(*Filter)

// FilterLevel with filter level.
func FilterLevel(level Level) FilterOption {
	return func(opts *Filter) {
		opts.level = level
	}
}

// FilterKey with filter key.
func FilterKey(key ...string) FilterOption {
	return func(o *Filter) {
		for _, v := range key {
			o.key[v] = struct{}{}
		}
	}
}

// FilterValue with filter value.
func FilterValue(value ...string) FilterOption {
	return func(o *Filter) {
		for _, v := range value {
			o.value[v] = struct{}{}
		}
	}
}

// Filter is a logger filter.
type Filter struct {
	logger Logger
	level  Level
	key    map[interface{}]struct{}
	value  map[interface{}]struct{}
}

// NewFilter new a logger filter.
func NewFilter(logger Logger, opts ...FilterOption) *Filter {
	options := Filter{
		logger: logger,
		key:    make(map[interface{}]struct{}),
		value:  make(map[interface{}]struct{}),
	}
	for _, o := range opts {
		o(&options)
	}
	return &options
}

func (f *Filter) Debug(msg string, fields ...zap.Field) {
	f.logger.Debug(msg, fields...)
}
func (f *Filter) Info(msg string, fields ...zap.Field) {
	f.logger.Info(msg, fields...)
}
func (f *Filter) Warn(msg string, fields ...zap.Field) {
	f.logger.Warn(msg, fields...)
}
func (f *Filter) Error(msg string, fields ...zap.Field) {
	f.logger.Error(msg, fields...)
}
func (f *Filter) DPanic(msg string, fields ...zap.Field) {
	f.logger.DPanic(msg, fields...)
}
func (f *Filter) Panic(msg string, fields ...zap.Field) {
	f.logger.Panic(msg, fields...)
}
func (f *Filter) Fatal(msg string, fields ...zap.Field) {
	f.logger.Fatal(msg, fields...)
}

func (f *Filter) formatFields(fields ...zap.Field) {
	if len(f.key) > 0 || len(f.value) > 0 {
		for i, field := range fields {
			if _, ok := f.key[field.Key]; ok {
				fields[i].Interface = nil
				fields[i].Integer = 0
				fields[i].String = "***"
			}
			if _, ok := f.value[field.String]; ok {
				fields[i].String = "***"
			}
		}
	}
}
