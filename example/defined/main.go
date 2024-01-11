package main

import (
	"github.com/yuanbaopig/logger"
	"go.uber.org/zap"
	"os"
)

func main() {

	opt := &logger.Options{
		OutputPath:    os.Stdout,
		Level:         "debug",
		DisableCaller: true,
	}

	newSlog := logger.Init(opt)
	defer newSlog.Sync()

	newSlog.Info("test")

	newSlog.SetOptions(logger.WithFields(zap.Int("userID", 10), zap.String("requestID", "fbf54504")))

	newSlog.Info("test1")
}
