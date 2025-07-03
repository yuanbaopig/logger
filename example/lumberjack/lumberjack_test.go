package main

import (
	"github.com/yuanbaopig/logger"
	"testing"
)

func TestLumberjack(t *testing.T) {

	log := logger.Log

	log.LumberjackLogger("test.log")

	log.Info("test")

}
