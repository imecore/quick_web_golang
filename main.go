package main

import (
	"github.com/joho/godotenv"
	"quick_web_golang/config"
	"quick_web_golang/engine"
	"quick_web_golang/log"
	"quick_web_golang/provider"
	"strconv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}

	logLevel, _ := strconv.Atoi(config.Get(config.LogLevel))
	log.SetLevel(logLevel)
	engine.Init()
	engine.Start()

	interrupt := provider.SignalWaitForInterrupt()
	_ = log.Infof("Captured %v, shutdown requested.\n", interrupt)
	engine.Stop()
	_ = log.Info("Exiting.")
}
