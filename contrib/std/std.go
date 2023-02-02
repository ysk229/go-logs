package std

import (
	"bytes"
	"fmt"
	"io"
	l "log"
	"sync"

	"github.com/mattn/go-colorable"

	log "github.com/ysk229/go-logs"
)

var _ log.Logger = (*stdLogger)(nil)

type stdLogger struct {
	log  *l.Logger
	pool *sync.Pool
}

// NewStdLogger new a logger with writer.
func NewStdLogger(w io.Writer) log.Logger {
	return &stdLogger{
		log: l.New(w, "", 0),
		pool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}
}

// Log print the kv pairs log.
func (l *stdLogger) Log(level log.Level, keyvals ...interface{}) error {
	if len(keyvals) == 0 {
		return nil
	}
	if (len(keyvals) & 1) == 1 {
		keyvals = append(keyvals, "KEYVALS UNPAIRED")
	}
	l.log.SetOutput(colorable.NewColorableStdout())
	buf := l.pool.Get().(*bytes.Buffer)
	var h, b string
	isMsg := false
	for i := 0; i < len(keyvals); i += 2 {
		if keyvals[i] == "ts" {
			h = ColorLog(level.String(), "[%s] %-8v", keyvals[i+1], level.String())
			continue
		}

		if keyvals[i] == "msg" {
			h += fmt.Sprintf("%-3v", keyvals[i+1])
			isMsg = true
			continue
		}
		b += fmt.Sprintf("%s=%v", ColorLog(level.String(), " %s", keyvals[i]), keyvals[i+1])
	}
	buf.WriteString(h)
	if !isMsg {
		buf.WriteString(fmt.Sprintf("%-4s", ""))
	}
	buf.WriteString(b)
	_ = l.log.Output(4, buf.String()) //nolint:gomnd
	buf.Reset()
	l.pool.Put(buf)
	return nil
}

func (l *stdLogger) SetLevel(level string) {
}

func (l *stdLogger) Close() error {
	return nil
}
