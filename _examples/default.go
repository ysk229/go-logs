package main

import (
	lg "log"

	"github.com/ysk229/go-logs/log"
)

func main() {
	l := log.DefaultLogger
	l.Info("info") // 不打印
	l.Warnw("key", "value")
	l.Errorf("%v", "sdfdsfds")
	l.SetLevel("info")
	lg.Println("test")
	l.Info("show info log")
	l.SetLevel("warn")
	l.Warnw("key", "value")
}
