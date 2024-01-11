package main

import (
	"github.com/yuanbaopig/logger"
)

func main() {
	logger.SetOptions(logger.WithLevel("debug"), logger.WithStdLevel("debug"), logger.WithDisableCaller(true))
	logger.SetOptions(logger.WithOutputPath(logger.GetFileLogWriter("test.log")))
	logger.SetOptions(logger.WithMultiCore(true))
	slog := logger.SLogger

	defer slog.Sync()

	slog.Debug("debug")
	slog.Info("info")
	logger.SetOptions(logger.WithFormatter("console"))
	slog.Warn("warning")
	slog.Error("error")

	// 自定义字段
	logger.SLogger.Infow("test", "k", "1")

}
