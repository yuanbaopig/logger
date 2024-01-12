package logger

import (
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// GetFileLogWriter 日志输出对象配置
func GetFileLogWriter(file string) zapcore.WriteSyncer {
	lj := &lumberjack.Logger{
		Filename:   file,  // 日志文件名
		MaxSize:    1024,  // 最大大小（MB）
		MaxBackups: 2,     // 最大备份数量
		MaxAge:     365,   // 最大保存天数
		Compress:   false, // 是否压缩
	}
	return zapcore.AddSync(lj)
}
