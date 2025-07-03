package options

import (
	"github.com/yuanbaopig/logger"
	"go.uber.org/zap"
	"testing"
)

func TestSetOption(t *testing.T) {
	lg := logger.New(
		logger.WithFields(zap.String("k", "v")),
		logger.WithLevel("info"),
	)
	lg.Debug("test")

	lg.SetOptions(logger.WithLevel("debug"))

	lg.Debug("test1")
}
