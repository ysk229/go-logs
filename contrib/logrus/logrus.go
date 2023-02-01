package logrus

import (
	"io"
	lg "log"

	"github.com/sirupsen/logrus"

	"github.com/ysk229/go-logs"
	"github.com/ysk229/go-logs/config"
	"github.com/ysk229/go-logs/file"
)

var _ log.Logger = (*Logger)(nil)

type Logger struct {
	logrus *logrus.Logger
}

func (l *Logger) SetLevel(level string) {
	// 日志级别
	parseLevel, e := logrus.ParseLevel(level)
	if e != nil {
		parseLevel = logrus.WarnLevel
	}
	l.logrus.SetLevel(parseLevel)
}

func (l *Logger) Log(level log.Level, keyvals ...interface{}) (err error) {
	var (
		logrusLevel logrus.Level
		fields      logrus.Fields = make(map[string]interface{})
		msg         string
	)

	switch level {
	case log.LevelDebug:
		logrusLevel = logrus.DebugLevel
	case log.LevelInfo:
		logrusLevel = logrus.InfoLevel
	case log.LevelWarn:
		logrusLevel = logrus.WarnLevel
	case log.LevelError:
		logrusLevel = logrus.ErrorLevel
	default:
		logrusLevel = logrus.DebugLevel
	}

	if logrusLevel > l.logrus.Level {
		return
	}

	if len(keyvals) == 0 {
		return nil
	}
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, "")
	}
	for i := 0; i < len(keyvals); i += 2 {
		key, ok := keyvals[i].(string)
		if !ok {
			continue
		}
		if key == logrus.FieldKeyMsg {
			msg, _ = keyvals[i+1].(string)
			continue
		}
		fields[key] = keyvals[i+1]
	}

	if len(fields) > 0 {
		l.logrus.WithFields(fields).Log(logrusLevel, msg)
	} else {
		l.logrus.Log(logrusLevel, msg)
	}

	return
}

func New(conf *config.Config) log.Logger {
	logger := logrus.New()
	// 日志级别
	parseLevel, e := logrus.ParseLevel(conf.Level)
	if e != nil {
		parseLevel = logrus.WarnLevel
	}
	logger.SetLevel(parseLevel)
	// log
	lg.SetFlags(lg.Lshortfile)
	lg.SetOutput(logger.Writer())
	getLog(conf, logger)
	return &Logger{
		logrus: logger,
	}
}

func getLog(conf *config.Config, logger *logrus.Logger) {
	var fileOpts []file.FileLogOption
	if len(conf.File.Path) > 0 {
		fileOpts = append(fileOpts, file.Path(conf.File.Path))
	}
	if len(conf.File.Mode) > 0 {
		fileOpts = append(fileOpts, file.Mode(conf.File.Mode))
	}
	if conf.File.MaxAge > 0 {
		fileOpts = append(fileOpts, file.MaxAge(conf.File.MaxAge))
	}
	if conf.File.Size > 0 {
		fileOpts = append(fileOpts, file.Size(conf.File.Size))
	}

	w := NewWrite(logger)
	switch conf.Both {
	case "all":
		// log
		w.setConsoleFormatter(conf.Format)
		w.writeFile(conf.Format, fileOpts...)
	case "file":
		w.writeFile(conf.Format, fileOpts...)
		w.logger.SetOutput(io.Discard)
	default:
		// log
		lg.SetOutput(logger.Writer())
		w.setConsoleFormatter(conf.Format)
	}
}
