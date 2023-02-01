package log

import (
	"fmt"
	"log"

	log2 "github.com/ysk229/go-logs"
	"github.com/ysk229/go-logs/config"
	"github.com/ysk229/go-logs/contrib/logrus"
	"github.com/ysk229/go-logs/contrib/std"
	"github.com/ysk229/go-logs/contrib/zap"
)

var (
	// DefaultMessageKey default message key.
	DefaultMessageKey = "msg"

	// DefaultLogger is default logger.
	DefaultLogger          = New(&config.Config{})
	_             log2.Log = (*l)(nil)
)

type l struct {
	log    log2.Logger
	msgKey string
}

// Option is WrapLogger option.
type Option func(*l)

// WithMessageKey with message key.
func WithMessageKey(k string) Option {
	return func(opts *l) {
		opts.msgKey = k
	}
}

func New(conf *config.Config, opts ...Option) log2.Log {
	logType := ""
	if conf != nil {
		logType = conf.Type // "zap" //std zap logrus
	}
	var logger log2.Logger
	switch logType {
	case "logrus":
		l := logrus.New(conf)
		logger = log2.With(l, "caller", log2.Caller(4), "type", "logrus")
	case "std":
		l := std.NewStdLogger(log.Writer())
		logger = log2.With(l, "ts", log2.DefaultTimestamp,
			"caller", log2.Caller(4), "type", "std")
	default:
		l := zap.New(conf)

		logger = log2.With(l, "caller", log2.Caller(4), "type", "zap")
	}
	optLog := &l{log: logger, msgKey: DefaultMessageKey}
	for _, o := range opts {
		o(optLog)
	}
	return optLog
}

func (l *l) SetLevel(level string) {
	l.log.SetLevel(level)
}

func (l *l) Log(level log2.Level, keyvals ...interface{}) error {
	return l.log.Log(level, keyvals...)
}

func (l *l) Info(a ...interface{}) {
	_ = l.log.Log(log2.LevelInfo, l.msgKey, fmt.Sprint(a...))
}

func (l *l) Warn(a ...interface{}) {
	_ = l.log.Log(log2.LevelWarn, l.msgKey, fmt.Sprint(a...))
}

func (l *l) Error(a ...interface{}) {
	_ = l.log.Log(log2.LevelError, l.msgKey, fmt.Sprint(a...))
}

func (l *l) Debug(a ...interface{}) {
	_ = l.log.Log(log2.LevelDebug, l.msgKey, fmt.Sprint(a...))
}

func (l *l) Debugf(format string, a ...interface{}) {
	_ = l.log.Log(log2.LevelDebug, l.msgKey, fmt.Sprintf(format, a...))
}

func (l *l) Debugw(keyvals ...interface{}) {
	_ = l.log.Log(log2.LevelDebug, keyvals...)
}

func (l *l) Infof(format string, a ...interface{}) {
	_ = l.log.Log(log2.LevelInfo, l.msgKey, fmt.Sprintf(format, a...))
}

func (l *l) Infow(keyvals ...interface{}) {
	_ = l.log.Log(log2.LevelInfo, keyvals...)
}

func (l *l) Warnf(format string, a ...interface{}) {
	_ = l.log.Log(log2.LevelWarn, l.msgKey, fmt.Sprintf(format, a...))
}

func (l *l) Warnw(keyvals ...interface{}) {
	_ = l.log.Log(log2.LevelWarn, keyvals...)
}

func (l *l) Errorf(format string, a ...interface{}) {
	_ = l.log.Log(log2.LevelError, l.msgKey, fmt.Sprintf(format, a...))
}

func (l *l) Errorw(keyvals ...interface{}) {
	_ = l.log.Log(log2.LevelError, keyvals...)
}

func (l *l) Fatal(a ...interface{}) {
	_ = l.log.Log(log2.LevelFatal, l.msgKey, fmt.Sprint(a...))
}

func (l *l) Fatalf(format string, a ...interface{}) {
	_ = l.log.Log(log2.LevelFatal, l.msgKey, fmt.Sprintf(format, a...))
}

func (l *l) Fatalw(keyvals ...interface{}) {
	_ = l.log.Log(log2.LevelFatal, keyvals...)
}

func (l *l) Print(v ...interface{}) {
	_ = l.log.Log(log2.LevelInfo, l.msgKey, fmt.Sprint(v...))
}

func (l *l) Println(v ...interface{}) {
	_ = l.log.Log(log2.LevelInfo, l.msgKey, fmt.Sprint(v...))
}

func (l *l) Fatalln(v ...interface{}) {
	_ = l.log.Log(log2.LevelFatal, l.msgKey, fmt.Sprint(v...))
}
