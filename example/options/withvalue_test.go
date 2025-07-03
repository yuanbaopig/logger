package options

import (
	"github.com/yuanbaopig/logger"
	"go.uber.org/zap"
	"testing"
)

func TestWithValue(t *testing.T) {

	logger.Log.Info("test")
	// 定义字段，生成新对象
	lv := logger.WithValues(zap.Int("userID", 10))
	lv.Info("test1")

	// 原结构不受影响
	logger.Log.Info("new test")
}
