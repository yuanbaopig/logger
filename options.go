package logger

import (
	"go.uber.org/zap"
)

// Options log config option.
type Options struct {
	Level             string
	DisableCaller     bool
	Fields            []zap.Field
	Format            string
	Development       bool
	DisableStacktrace bool
	OutputPaths       []string
	ErrorOutputPaths  []string
	EnableColor       bool
	Name              string
}

type Option func(*Options)

// WithDisableStacktrace change log Stacktrace status.
func WithDisableStacktrace(enable bool) Option {
	return func(o *Options) {
		o.DisableStacktrace = enable
	}
}

// WithDevelopment change log development status.
func WithDevelopment(enable bool) Option {
	return func(o *Options) {
		o.Development = enable
	}
}

// WithLevel change log level, default info.
func WithLevel(level string) Option {
	return func(o *Options) {
		o.Level = level
	}
}

// WithEnableColor change log enable color.
func WithEnableColor(enable bool) Option {
	return func(o *Options) {
		o.EnableColor = enable
	}
}

// WithFields append logger fields.
func WithFields(fields ...zap.Field) Option {
	//zap.Fields() math return zapcore.Field
	//zap.Fields()
	return func(o *Options) {
		o.Fields = append(o.Fields, fields...)
	}

}

// WithFormat log format.
func WithFormat(format string) Option {
	return func(o *Options) {
		o.Format = format
	}
}

// WithOutputPaths change output object.
func WithOutputPaths(outputPaths []string) Option {
	return func(o *Options) {
		o.OutputPaths = outputPaths
	}
}

// WithErrorOutputPaths change error output object.
func WithErrorOutputPaths(errorOutputPaths []string) Option {
	/*
		ErrorOutputPaths 在 zap 配置中用于指定错误输出的路径。
		它主要用于记录 zap 日志库内部产生的错误，例如编码错误或写入日志时的失败。
		这与应用程序使用 logger.Error(...) 记录的错误日志不同。ErrorOutputPaths 更关注于日志系统内部的错误，而不是应用程序逻辑中产生的错误。
	*/
	return func(o *Options) {
		o.ErrorOutputPaths = errorOutputPaths
	}
}

// WithDisableCaller manage caller status.
func WithDisableCaller(caller bool) Option {
	return func(o *Options) {
		o.DisableCaller = caller
	}
}

// WithValues creates a child logger and adds Zap Fields to it.
func WithValues(keysAndValues ...zap.Field) *zap.Logger {
	return SLogger.WithValues(keysAndValues...)
}

func (l *log) WithValues(keysAndValues ...zap.Field) *zap.Logger {
	//newLogger := l.Sugar().With(zap.Fields(keysAndValues...))

	if len(keysAndValues) == 0 {
		return l.Logger
	}

	return l.With(keysAndValues...)

}
