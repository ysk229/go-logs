package main

import (
	"github.com/ysk229/go-logs/config"
	"github.com/ysk229/go-logs/log"
)

func main() {
	l := log.New(&config.Config{
		Type:  "logrus",
		Level: "info",
	})

	l.Info("info2222")
	l.Warnw("key", "value")
	l.Errorf("%v", "sdfdsfds")
	l.SetLevel("info")
	l.Info("show info log")
	l.SetLevel("warn")
	l.Warnw("key", "value")
}
