package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
	"time"
)

const (
	consoleFormat = "console"
	jsonFormat    = "json"
)

// SLogger global logger
var SLogger = New()

// New default build logger.
func New(opts ...Option) *log {
	// Options init

	l := &log{
		opt: &Options{
			DisableStacktrace: true,
			OutputPaths:       []string{"stdout"},
			ErrorOutputPaths:  []string{"stderr"},
		},
	}
	l.initOptions(opts...)
	// new sugaredLogger
	if err := l.newSugaredLogger(); err != nil {
		panic(err)
	}

	return l
}

// Init user-defined options to build logger.
func Init(opt *Options) *log {
	l := &log{opt: opt}
	if err := l.newSugaredLogger(); err != nil {
		panic(err)
	}
	return l
}

// getEncoderConfig for create encoder config.
func getEncoderConfig(EnableColor bool) zapcore.EncoderConfig {
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

	//encoderConfig := zap.NewProductionEncoderConfig()
	//encoderConfig.TimeKey = "ts" // 放在第一个位置
	//encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	//encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	encodeLevel := zapcore.CapitalLevelEncoder
	if EnableColor {
		encodeLevel = zapcore.CapitalColorLevelEncoder
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "ts",
		MessageKey:    "message",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   encodeLevel,
		//EncodeTime:     timeEncoder,		// 日期格式
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: milliSecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
	return encoderConfig
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func milliSecondsDurationEncoder(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendFloat64(float64(d) / float64(time.Millisecond))
}

// newSugaredLogger SugaredLogger init.
func (l *log) newSugaredLogger() error {
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

	// 是否开启编辑器颜色配置
	var EnableColor bool
	if l.opt.Format == consoleFormat && l.opt.EnableColor {
		EnableColor = true
	}

	//	暂存起来，别的功能用的到
	l.encoderConfig = getEncoderConfig(EnableColor)

	// 根据字符串解析，如果有报错则会默认info级别
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(l.opt.Level)); err != nil {
		fmt.Println(err)
		zapLevel = zapcore.InfoLevel
	}

	// 日志格式
	encoderFormat := getEncoderFormat(l.opt.Format)

	// zap config
	zc := &zap.Config{
		Level:             zap.NewAtomicLevelAt(zapLevel),
		Development:       l.opt.Development,
		DisableCaller:     l.opt.DisableCaller,
		DisableStacktrace: l.opt.DisableStacktrace,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         encoderFormat,
		EncoderConfig:    l.encoderConfig,
		OutputPaths:      l.opt.OutputPaths,
		ErrorOutputPaths: l.opt.ErrorOutputPaths,
	}

	// AddStacktrace 控制在哪个级别会输出Stacktrace，此处设置为panic级别
	//logger, err := zc.Build(zap.AddStacktrace(zapcore.PanicLevel))

	logger, err := zc.Build(zap.Fields(l.opt.Fields...))
	if err != nil {
		return err
	}
	// 将标准库的 log 输出重定向到 zap
	//zap.RedirectStdLog(logger.Named(l.opt.Name))

	/*
		当你使用 `zap.ReplaceGlobals(newLogger)` 设置全局日志记录器时，
		`zap.L()` 将返回你最后一次调用 `zap.ReplaceGlobals` 时所设置的 `newLogger` 实例。
		在 `zap` 中，全局日志记录器是唯一的，意味着你不能同时拥有多个全局日志记录器实例。

		如果你多次调用 `zap.ReplaceGlobals`，每次传入一个不同的 `logger.Named(name)` 实例，
		`zap.L()` 将总是引用最后一次设置的那个日志记录器实例。前面设置的日志记录器将被覆盖，不再是全局日志记录器。
	*/
	zap.ReplaceGlobals(logger)
	l.Logger = zap.L()
	return nil
}

// getEncoderFormat create encoder based on the format as the foundation.
func getEncoderFormat(format string) string {
	var encoder string

	switch format {
	case "json":
		// console 日志格式
		encoder = jsonFormat
	case "console":
		// json 日志格式
		encoder = consoleFormat
	default:
		encoder = jsonFormat
	}
	return encoder
}

type log struct {
	*zap.Logger
	opt           *Options
	mu            sync.Mutex
	encoderConfig zapcore.EncoderConfig
}

// initOptions init config option of log.
func (l *log) initOptions(opts ...Option) {
	for _, opt := range opts {
		opt(l.opt)
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
	if err := l.newSugaredLogger(); err != nil {
		panic(err)
	}
}

// WithName adds a new path segment to the logger's name. Segments are joined by
// periods. By default, Loggers are unnamed.
func WithName(s string) *zap.Logger {
	return SLogger.Named(s)
}

func (l *log) LumberjackLogger(filename string) {

	// 创建 lumberjack.Logger 实例作为日志输出
	lumberjackLogger := GetFileLogWriter(filename)

	// 创建新的 WriteSyncer
	writeSyncer := zapcore.AddSync(lumberjackLogger)

	// 文件中输出有颜色会变成乱码，所以强制修改
	l.encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	var encoder zapcore.Encoder
	switch l.opt.Format {
	case jsonFormat:
		encoder = zapcore.NewJSONEncoder(l.encoderConfig)
	case consoleFormat:
		encoder = zapcore.NewConsoleEncoder(l.encoderConfig)
	default:
		encoder = zapcore.NewJSONEncoder(l.encoderConfig)
	}

	// 使用原有的 EncoderConfig 创建新的 zapcore.Core
	newCore := zapcore.NewCore(
		encoder,
		writeSyncer,
		l.Level(),
	)

	// 更新 logger
	logger := l.Logger.WithOptions(zap.WrapCore(func(zapcore.Core) zapcore.Core {
		return newCore
	}))

	// 更新全局日志记录器
	zap.ReplaceGlobals(logger)
	l.Logger = zap.L()
}
