package std

import (
	l "log"
	"testing"

	"github.com/ysk229/go-logs"
)

func TestStdLogger(t *testing.T) {
	logger := NewStdLogger(l.Writer())
	logger = log.With(logger, "caller", log.DefaultCaller, "ts", log.DefaultTimestamp)

	_ = logger.Log(log.LevelInfo, "msg", "test debug")

	_ = logger.Log(log.LevelInfo, "msg", "test info")
	_ = logger.Log(log.LevelInfo, "msg", "test warn")
	_ = logger.Log(log.LevelInfo, "msg", "test error")
	_ = logger.Log(log.LevelDebug, "singular", "test")
	_ = logger.Log(log.LevelWarn, "msg", "test error")
	_ = logger.Log(log.LevelWarn, "warn singular", "sdfsdfsdf")
}
