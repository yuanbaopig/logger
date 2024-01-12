package main

import (
	"github.com/yuanbaopig/logger"
	"go.uber.org/zap"
)

func main() {

	logger.SetOptions(
		logger.WithFormat("console"),                                         // 日志格式
		logger.WithEnableColor(true),                                         // 是否显示日志级别颜色，只有在console下才支持
		logger.WithDisableCaller(false),                                      // 是否关闭调用信息
		logger.WithLevel("debug"),                                            // 日志级别
		logger.WithDisableStacktrace(true),                                   // 是否关闭Stacktrace
		logger.WithFields(zap.String("app", "myApp"), zap.Int("version", 1)), // 日志结构化
		//logger.WithDevelopment(true),
		//logger.WithOutputPaths([]string{"path_test.log"}),
	)

	//lg := logger.SLogger.Sugar()
	lg := zap.L().Sugar()
	lg.Debug("debug")
	lg.Info("info")
	lg.Warn("warning")
	lg.Error("error")

	zap.L().Info("test", zap.Int("key", 1))
	lg.Infow("test", "X-Request-ID", "fbf54504-64da-4088-9b86-67824a7fb508")
	lg.Fatal("fatal")
}
