package log

type BaseLog interface {
	Logger
	Info(a ...interface{})
	Warn(a ...interface{})
	Error(a ...interface{})
}
type Log interface {
	BaseLog
	Debug(a ...interface{})
	Debugf(format string, a ...interface{})
	Debugw(keyvals ...interface{})

	Infof(format string, a ...interface{})
	Infow(keyvals ...interface{})

	Warnf(format string, a ...interface{})
	Warnw(keyvals ...interface{})

	Errorf(format string, a ...interface{})
	Errorw(keyvals ...interface{})

	Fatal(a ...interface{})
	Fatalf(format string, a ...interface{})
	Fatalw(keyvals ...interface{})
}
