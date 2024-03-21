package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func main() {
	// 创建一个 EncodingConfig。在这里可以配置时间格式、日志级别显示格式等。
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	errorLogFile, _ := os.Create("error.log")
	normalLogFile, _ := os.Create("normal.log")
	// 设置日志级别
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	// 创建不同的输出目的地
	errorWriter := zapcore.AddSync(errorLogFile)                                                          // 错误级别日志的输出目的地，这里设置为标准错误
	infoWriter := zapcore.NewMultiWriteSyncer(zapcore.AddSync(normalLogFile), zapcore.AddSync(os.Stdout)) // 信息和调试级别日志的输出目的地，这里设置为标准输出

	// 创建两个 Core，分别负责不同日志级别的日志。
	infoCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg), // 这里使用 JSON 编码器
		infoWriter,
		infoLevel,
	)
	errorCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg), // 这里使用 JSON 编码器
		errorWriter,
		errorLevel,
	)

	// 使用 Tee 来合并两个 Core
	core := zapcore.NewTee(infoCore, errorCore)

	// 创建 Logger
	logger := zap.New(core, zap.AddStacktrace(zap.InfoLevel))

	defer logger.Sync() // 将缓冲中的日志进行同步

	// 使用示例
	logger.Info("这是一条信息级别的日志")
	logger.Debug("这是一条调试级别的日志")
	logger.Error("这是一条错误级别的日志")
}
