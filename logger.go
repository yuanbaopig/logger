package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
	"time"
)

const (
	consoleFormat = "console"
	jsonFormat    = "json"
)

// Log global logger
var Log = New()

// New default build logger.
func New(opts ...Option) *log {
	// Options init
	l := &log{
		opt: &Options{
			DisableStacktrace: true,
			OutputPaths:       []string{"stdout"},
		},
	}
	l.initOptions(opts...)
	// new sugaredLogger
	if err := l.newLogger(); err != nil {
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
		MessageKey:    "msg",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   encodeLevel,
		//EncodeTime:     timeEncoder,		// 日期格式
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: milliSecondsDurationEncoder,
		//EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeName:   zapcore.FullNameEncoder,
	}
	return encoderConfig
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func milliSecondsDurationEncoder(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendFloat64(float64(d) / float64(time.Millisecond))
}

func createWriteSyncersAndCore(paths []string, encoder zapcore.Encoder, levelEnabler zapcore.LevelEnabler) (zapcore.Core, error) {
	/*
		var normalOutputList []zapcore.WriteSyncer
		for _, fileName := range l.opt.OutputPaths {
			ws, err := openFileWriteSyncer(fileName)
			if err != nil {
				return err
			}
			normalOutputList = append(normalOutputList, ws)
		}

		var normalCore zapcore.Core
		normalCore = zapcore.NewCore(
			encoder,
			zapcore.NewMultiWriteSyncer(normalOutputList...),
			infoLevel,
		)


		var errorCore zapcore.Core

		if len(l.opt.ErrorOutputPaths) > 0 {
			// error log level for output
			errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl >= zapcore.ErrorLevel
			})

			var errorOutputList []zapcore.WriteSyncer
			for _, fileName := range l.opt.ErrorOutputPaths {
				ws, err := openFileWriteSyncer(fileName)
				if err != nil {
					return err
				}
				errorOutputList = append(errorOutputList, ws)
			}

			errorCore = zapcore.NewCore(
				encoder,
				zapcore.NewMultiWriteSyncer(errorOutputList...),
				errorLevel,
			)
		}

	*/
	var writeSyncers []zapcore.WriteSyncer
	for _, fileName := range paths {
		ws, err := openFileWriteSyncer(fileName)
		if err != nil {
			return nil, err
		}
		writeSyncers = append(writeSyncers, ws)
	}
	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(writeSyncers...),
		levelEnabler,
	)
	return core, nil
}

func openFileWriteSyncer(fileName string) (zapcore.WriteSyncer, error) {
	if fileName == "stdout" {
		return zapcore.AddSync(os.Stdout), nil
	}
	if fileName == "stderr" {
		return zapcore.AddSync(os.Stderr), nil
	}

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return zapcore.AddSync(file), nil
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

// log logger options
type log struct {
	*zap.Logger
	opt           *Options
	mu            sync.Mutex
	encoderConfig zapcore.EncoderConfig
}

// newLogger constructs a Logger.
func (l *log) newLogger() error {
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
		//fmt.Println(err)
		zapLevel = zapcore.InfoLevel
	}

	// 日志格式
	encoderFormat := getEncoderFormat(l.opt.Format)
	//
	//// zap config
	//zc := &zap.Config{
	//	Level:             zap.NewAtomicLevelAt(zapLevel),
	//	Development:       l.opt.Development,
	//	DisableCaller:     l.opt.DisableCaller,
	//	DisableStacktrace: l.opt.DisableStacktrace,
	//	Sampling: &zap.SamplingConfig{
	//		Initial:    100,
	//		Thereafter: 100,
	//	},
	//	Encoding:         encoderFormat,
	//	EncoderConfig:    l.encoderConfig,
	//	OutputPaths:      l.opt.OutputPaths,
	//	ErrorOutputPaths: l.opt.ErrorOutputPaths,
	//}
	//
	//// AddStacktrace 控制在哪个级别会输出Stacktrace，此处设置为panic级别
	////logger, err := zc.Build(zap.AddStacktrace(zapcore.PanicLevel))
	//
	//logger, err := zc.Build(zap.Fields(l.opt.Fields...))
	//if err != nil {
	//	return err
	//}

	// new version
	// --------------------------------------

	var encoder zapcore.Encoder
	switch encoderFormat {
	case jsonFormat:
		encoder = zapcore.NewJSONEncoder(l.encoderConfig)
	case consoleFormat:
		encoder = zapcore.NewConsoleEncoder(l.encoderConfig)
	}

	var coreList []zapcore.Core
	{
		// general log level for output
		infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapLevel
		})
		// 在使用的地方生成 normalCore
		normalCore, err := createWriteSyncersAndCore(l.opt.OutputPaths, encoder, infoLevel)
		if err != nil {
			return err
		}
		coreList = append(coreList, normalCore)
	}

	// 根据条件，生成 errorCore
	if len(l.opt.ErrorOutputPaths) > 0 {
		errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel
		})

		errorCore, err := createWriteSyncersAndCore(l.opt.ErrorOutputPaths, encoder, errorLevel)
		if err != nil {
			return err
		}

		coreList = append(coreList, errorCore)

	}

	core := zapcore.NewTee(coreList...)

	var zapOpts []zap.Option
	if !l.opt.DisableCaller { // false
		zapOpts = append(zapOpts, zap.AddCaller())
	}

	if !l.opt.DisableStacktrace { // false
		// AddStacktrace 控制在哪个级别会输出Stacktrace，此处设置为error级别
		zapOpts = append(zapOpts, zap.AddStacktrace(zap.ErrorLevel))
	}

	if len(l.opt.Fields) != 0 {
		zapOpts = append(zapOpts, zap.Fields(l.opt.Fields...))
	}

	if l.opt.CallerSkip != 0 {
		zapOpts = append(zapOpts, zap.AddCallerSkip(l.opt.CallerSkip))
	}

	logger := zap.New(core, zapOpts...)

	// --------------------------------------
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
	l.Logger = logger
	return nil
}

// initOptions init config option of log.
func (l *log) initOptions(opts ...Option) {
	for _, opt := range opts {
		opt(l.opt)
	}
}

// SetOptions user reset defined.
func (l *log) SetOptions(opts ...Option) {
	l.mu.Lock()
	defer l.mu.Unlock()
	// change Options
	l.initOptions(opts...)
	if err := l.newLogger(); err != nil {
		panic(err)
	}
}

// LumberjackLogger create a lumberjack log file.
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
	l.Logger = logger
}
