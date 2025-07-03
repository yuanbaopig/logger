package options

import (
	"github.com/yuanbaopig/logger"
	"go.uber.org/zap"
	"testing"
)

func TestOptions(t *testing.T) {
	lg := logger.New(
		logger.WithFormat("console"),                                         // 日志格式
		logger.WithEnableColor(true),                                         // 是否显示日志级别颜色，只有在console下才支持
		logger.WithDisableCaller(false),                                      // 是否关闭调用信息，默认为false
		logger.WithLevel("info"),                                             // 日志级别
		logger.WithDisableStacktrace(false),                                  // 是否关闭Stacktrace，默认false
		logger.WithFields(zap.String("app", "myApp"), zap.Int("version", 1)), // 日志结构化
		//logger.WithOutputPaths([]string{"test.log"}),
		logger.WithOutputPaths([]string{"stdout"}),
		//logger.WithErrorOutputPaths([]string{"stderr", "error.log"}), // 如果使用同时使用stdout和stderr，并且没有屏蔽stdout的情况下会导致error信息输出两次
	)

	lg.Info("test")
	lg.Error("error")

}
