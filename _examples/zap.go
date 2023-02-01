package main

import (
	"fmt"

	"github.com/ysk229/go-logs/config"
	"github.com/ysk229/go-logs/log"
)

func main() {
	l := log.New(&config.Config{
		Type:   "zap",
		Level:  "info",
		Format: "json",
	})

	l.Info("info2222")
	l.Warnw("key", "value")
	l.Errorf("%v", "sdfdsfds")
	l.SetLevel("info")
	l.Info("show info log")
	l.SetLevel("warn")
	l.Warnw("key", "value")
	l = log.New(&config.Config{
		Both:   "all",
		Type:   "zap",
		Level:  "info",
		Format: "text",
		File:   config.FileConf{Path: "./../logs/"},
	})

	l.Info("info2222")
	l.Warnw("key", "value")
	l.Errorf("%v", "sdfdsfds")
	l.SetLevel("info")
	l.Info("show info log")
	l.SetLevel("warn")
	l.Warnw("key", "value")
	fmt.Println("sdfsdfsdfs")
}
