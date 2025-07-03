package main

import (
	"context"
	"github.com/yuanbaopig/logger"
	"testing"
)

func TestContext(t *testing.T) {
	defer logger.Log.Sync()
	
	lv := logger.Log.Logger

	lv.Info("test")

	// 讲logger存储到context中
	ctx := logger.WithContext(context.Background(), lv)

	// 进行context传递
	PrintString(ctx, "World")
}

func PrintString(ctx context.Context, str string) {
	// 从context中获取logger
	lc := logger.FromContext(ctx)
	lc.Sugar().Infof("Hello %s", str)
}
