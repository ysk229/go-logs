package log

import (
	log2 "log"
	"testing"

	"github.com/ysk229/go-logs/config"
)

func TestUseZapLog(t *testing.T) {
	l := New(&config.Config{
		Type:  "zap",
		Both:  "all",
		Level: "info",
		// Format: "json",
	})
	log2.Println("23423423423sd")
	l.Info("test.....")
	l.Warnw("sdfdsfds", "sdfsdfsd")
	// l.Errorf("%s test", "324234234")
}

func TestUseLogrusLog(t *testing.T) {
	l := New(&config.Config{
		Type:  "logrus",
		Both:  "all",
		Level: "info",
		// Format: "json",
	})
	log2.Println("23423423423sd")
	l.Info("test.....")
	l.Warnw("sdfdsfds", "sdfsdfsd")
	// l.Errorf("%s test", "324234234")
}

func TestUseStdLog(t *testing.T) {
	l := New(&config.Config{
		Type:  "std",
		Both:  "all",
		Level: "info",
		// Format: "json",
	})
	log2.Println("23423423423sd")
	l.Info("test.....")
	l.Info("test info.....")
	l.Info("test info2.....")
	l.Warnw("sdfdsfds", "sdfsdfsd")
	// l.Errorf("%s test", "324234234")
}

func TestLog(t *testing.T) {
	l := DefaultLogger

	l.Info("info") // 不打印
	l.Warnw("key", "value")
	l.Errorf("%v", "sdfdsfds")
	l.SetLevel("info")
	l.Info("show info log")
	l.SetLevel("warn")
	l.Warnw("key", "value")
}
