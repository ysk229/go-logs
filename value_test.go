package log_test

import (
	"context"
	"log"
	"testing"

	l "github.com/ysk229/go-logs"
	"github.com/ysk229/go-logs/contrib/std"
)

func TestValue(t *testing.T) {
	logger := std.NewStdLogger(log.Writer())
	logger = l.With(logger, "ts", l.DefaultTimestamp, "caller", l.DefaultCaller)
	_ = logger.Log(l.LevelInfo, "msg", "helloworld")

	logger = std.NewStdLogger(log.Writer())
	logger = l.With(logger)
	_ = logger.Log(l.LevelDebug, "msg", "helloworld")

	var v1 interface{}
	got := l.Value(context.Background(), v1)
	if got != v1 {
		t.Errorf("Value() = %v, want %v", got, v1)
	}
	var v2 l.Valuer = func(ctx context.Context) interface{} {
		return 3
	}
	got = l.Value(context.Background(), v2)
	res := got.(int)
	if res != 3 {
		t.Errorf("Value() = %v, want %v", res, 3)
	}
}
