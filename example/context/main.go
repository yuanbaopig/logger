package main

import (
	"context"
	"github.com/yuanbaopig/logger"
	"go.uber.org/zap"
)

func main() {
	defer logger.SLogger.Sync()

	// 定义字段
	lv := logger.WithValues(zap.Int("userID", 10))

	lv.Info("test")

	// 讲logger存储到context中
	ctx := logger.WithContext(context.Background(), lv)

	// 进行context传递
	PrintString(ctx, "World")

	// 原结构不受影响
	logger.SLogger.Sugar().Infof("Hello World")
}

func PrintString(ctx context.Context, str string) {
	//从context中获取logger
	lc := logger.FromContext(ctx)
	lc.Sugar().Infof("Hello %s", str)
}
