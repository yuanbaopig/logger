package main

import (
	"github.com/yuanbaopig/logger"
	"testing"
)

func TestSimple(t *testing.T) {

	// 调用包里的logger
	logger.Log.Info("info")
	defer logger.Log.Sync()

	// 调用SugaredLogger
	//slg := lg.Sugar()
	//slg.Infow("test", "X-Request-ID", "fbf54504-64da-4088-9b86-67824a7fb508")
	lg := logger.Log.Sugar()
	lg.Debug("debug")
	lg.Warn("warning")
	lg.Error("error")
	lg.Fatal("fatal")

}
