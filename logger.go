package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
)

// SLogger global logger
var SLogger = New()

// New default build logger.
func New(opts ...Option) *log {
	// Options init
	l := &log{
		opt: &Options{},
	}
	l.initOptions(opts...)
	// new sugaredLogger
	l.newSugaredLogger()
	return l
}

// Init user-defined options to build logger.
func Init(opt *Options) *log {
	l := &log{opt: opt}
	l.newSugaredLogger()
	return l
}

// getEncoderConfig for create encoder config.
func getEncoderConfig() zapcore.EncoderConfig {
	/*
		encoderConfig := zapcore.EncoderConfig{
			MessageKey:     "message",
			LevelKey:       "level",
			TimeKey:        "timestamp",
			NameKey:        "logger",
			CallerKey:      "caller",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    encodeLevel,
			EncodeTime:     timeEncoder,
			EncodeDuration: milliSecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}

	*/

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "ts" // 放在第一个位置
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	return encoderConfig
}

// newSugaredLogger SugaredLogger init.
func (l *log) newSugaredLogger() {

	/*
		loggerConfig := &zap.Config{
			Level:             zap.NewAtomicLevelAt(zapLevel),
			Development:       opts.Development,
			DisableCaller:     opts.DisableCaller,
			DisableStacktrace: opts.DisableStacktrace,
			Sampling: &zap.SamplingConfig{
				Initial:    100,
				Thereafter: 100,
			},
			Encoding:         opts.Format,
			EncoderConfig:    encoderConfig,
			OutputPaths:      opts.OutputPaths,
			ErrorOutputPaths: opts.ErrorOutputPaths,
		}
		l, err := loggerConfig.Build(zap.AddStacktrace(zapcore.PanicLevel), zap.AddCallerSkip(1))
	*/

	// 编辑器（配置）
	encoderConfig := getEncoderConfig()

	// 实例化AtomicLevel对象
	//Level := zap.NewAtomicLevel()
	//Level.SetLevel(l.opt.Level) // zap.DebugLevel 日志级别

	// 日志输出对象
	//l.opt.OutputPath := os.Stdout
	// 根据字符串解析，如果有报错则会默认info级别
	parsedLevel, _ := zapcore.ParseLevel(l.opt.Level)

	// 实例化 file core对象
	fileCore := zapcore.NewCore(getEncoder(encoderConfig, l.opt.Format), l.opt.OutputPath, parsedLevel)

	// stand out core
	var combinedCore zapcore.Core
	if l.opt.StdOutput {
		// 解析日志级别
		stdParsedLevel, _ := zapcore.ParseLevel(l.opt.StdLevel)
		// 获取一个开发模式的encoder config
		consoleEncoder := zap.NewDevelopmentEncoderConfig()
		// new console core
		consoleCore := zapcore.NewCore(zapcore.NewConsoleEncoder(consoleEncoder), os.Stdout, stdParsedLevel)
		// 将 fileCore 和 consoleCore 组合成一个输出
		combinedCore = zapcore.NewTee(fileCore, consoleCore)
	} else {
		combinedCore = fileCore
	}

	// new logger
	logger := zap.New(combinedCore, zap.WithCaller(l.opt.DisableCaller), zap.Fields(l.opt.Fields...))
	// get SugaredLogger
	l.SugaredLogger = logger.Sugar()

}

// getEncoder create encoder based on the format as the foundation.
func getEncoder(config zapcore.EncoderConfig, format string) zapcore.Encoder {
	var encoder zapcore.Encoder

	switch format {
	case "json":
		// console 日志格式
		encoder = zapcore.NewJSONEncoder(config)
	case "console":
		// json 日志格式
		encoder = zapcore.NewConsoleEncoder(config)
	default:
		encoder = zapcore.NewJSONEncoder(config)
	}
	return encoder
}

type log struct {
	*zap.SugaredLogger
	opt *Options
	mu  sync.Mutex
}

// initOptions init config option of log.
func (l *log) initOptions(opts ...Option) {
	for _, opt := range opts {
		opt(l.opt)
	}

	// if OutputPath is nil then set default value.
	if l.opt.OutputPath == nil {
		l.opt.OutputPath = os.Stderr
	}
}

// SetOptions global sugared logger use.
func SetOptions(opts ...Option) {
	SLogger.SetOptions(opts...)
}

// SetOptions user-defined sugared logger use.
func (l *log) SetOptions(opts ...Option) {
	l.mu.Lock()
	defer l.mu.Unlock()
	// change Options
	l.initOptions(opts...)
	// reset sugared logger
	l.newSugaredLogger()
}

// WithName adds a new path segment to the logger's name. Segments are joined by
// periods. By default, Loggers are unnamed.
func WithName(s string) *log { return SLogger.WithName(s) }

func (l *log) WithName(name string) *log {
	newLogger := l.SugaredLogger.Named(name)
	l.SugaredLogger = newLogger
	return l
}
