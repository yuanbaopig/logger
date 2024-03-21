package main

import (
	"github.com/yuanbaopig/logger"
	"go.uber.org/zap"
)

func main() {

	logger.SetOptions(
		logger.WithFormat("json"),                                            // 日志格式
		logger.WithEnableColor(true),                                         // 是否显示日志级别颜色，只有在console下才支持
		logger.WithDisableCaller(false),                                      // 是否关闭调用信息
		logger.WithLevel("info"),                                             // 日志级别
		logger.WithDisableStacktrace(true),                                   // 是否关闭Stacktrace
		logger.WithFields(zap.String("app", "myApp"), zap.Int("version", 1)), // 日志结构化
		logger.WithOutputPaths([]string{"stdout", "test.log"}),
		//logger.WithErrorOutputPaths([]string{"stderr","error.log"}),
	)

	//直接调用全局zap logger
	zap.L().Info("test", zap.Int("key", 1))
	// 调用包里的logger
	lg := logger.Log
	defer lg.Sync()
	lg.Info("info")

	// 调用SugaredLogger
	slg := lg.Sugar()
	slg.Infow("test", "X-Request-ID", "fbf54504-64da-4088-9b86-67824a7fb508")

	lg.Debug("debug")
	lg.Warn("warning")
	lg.Error("error")
	lg.Fatal("fatal")

}
