package hook

import (
	"github.com/sirupsen/logrus"
	"runtime"
)

type JsonFormatter struct {
	*logrus.JSONFormatter
}

//const TimestampFormat = "2006-01-02T15:04:05.999Z07:00"

const TimestampFormat = "2006-01-02.15:04:05.000000"

func NewJsonFormatter() logrus.Formatter {
	return &JsonFormatter{
		JSONFormatter: &logrus.JSONFormatter{
			TimestampFormat: TimestampFormat,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "ts",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "msg",
			},
			CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
				return "", ""
			},
		}}
}

func (f *JsonFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return f.JSONFormatter.Format(entry)
}
