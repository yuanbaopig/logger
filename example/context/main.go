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
	// 讲logger存储到context中
	ctx := lv.WithContext(context.Background())

	// 进行context传递
	PrintString(ctx, "World")

	// 原结构不受影响
	logger.SLogger.Infof("Hello World")
}

func PrintString(ctx context.Context, str string) {
	// 从context中获取logger
	lc := logger.FromContext(ctx)
	lc.Infof("Hello %s", str)
}
