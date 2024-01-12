package main

import "github.com/yuanbaopig/logger"

func main() {

	log := logger.Init(&logger.Options{
		Level:       "debug",
		Format:      "json",
		OutputPaths: []string{"stdout"},
	})

	log.LumberjackLogger("test.log")

	log.Info("test")

	//logger.WithFields()
}
