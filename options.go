package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Options log config option.
type Options struct {
	OutputPath    zapcore.WriteSyncer
	Level         string
	StdLevel      string
	StdOutput     bool
	DisableCaller bool
	Fields        []zap.Field
	Format        string
}

type Option func(*Options)

// WithLevel change log level, default info.
func WithLevel(level string) Option {
	return func(o *Options) {
		o.Level = level
	}
}

// WithStdLevel change stand out log level, default info.
func WithStdLevel(level string) Option {
	return func(o *Options) {
		o.StdLevel = level
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

// WithFormatter log format
func WithFormatter(formatter string) Option {
	return func(o *Options) {
		o.Format = formatter
	}
}

// WithOutputPath change output object.
func WithOutputPath(output zapcore.WriteSyncer) Option {
	return func(o *Options) {
		o.OutputPath = output
	}
}

// WithDisableCaller manage caller status.
func WithDisableCaller(caller bool) Option {
	return func(o *Options) {
		o.DisableCaller = caller
	}
}

// WithValues creates a child logger and adds Zap Fields to it.
func WithValues(keysAndValues ...zap.Field) *log { return SLogger.WithValues(keysAndValues...) }

func (l *log) WithValues(keysAndValues ...zap.Field) *log {
	//newLogger := l.SugaredLogger.With(zap.Fields(keysAndValues...))

	newLog := &log{opt: l.opt}
	newLog.opt.Fields = keysAndValues

	newLog.newSugaredLogger()
	return newLog
}

// WithMultiCore use multi output, file and standard out
func WithMultiCore(std bool) Option {
	return func(o *Options) {
		o.StdOutput = std
	}
}
