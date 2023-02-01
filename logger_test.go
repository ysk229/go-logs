package log_test

import (
	"context"
	"log"
	"testing"

	l "github.com/ysk229/go-logs"
	"github.com/ysk229/go-logs/contrib/std"
)

func TestInfo(t *testing.T) {
	logger := std.NewStdLogger(log.Writer())
	logger = l.With(logger, "ts", l.DefaultTimestamp)
	logger = l.With(logger, "caller", l.DefaultCaller)

	_ = logger.Log(l.LevelInfo, "key1", "value1")
}

func TestWithContext(t *testing.T) {
	l.WithContext(context.Background(), nil)
}
