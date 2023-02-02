package logrus

import (
	"github.com/mattn/go-colorable"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"

	log "github.com/ysk229/go-logs"
	"github.com/ysk229/go-logs/contrib/logrus/hook"
	"github.com/ysk229/go-logs/file"
)

const Mode = "json"

type Write struct {
	logger *logrus.Logger
}

func NewWrite(logger *logrus.Logger) *Write {
	return &Write{logger: logger}
}

func (w *Write) setConsoleFormatter(encodeName string) {
	if encodeName == Mode {
		w.logger.SetFormatter(hook.NewJSONFormatter())
		w.logger.SetOutput(log.NewJSONColorable())
	} else {
		w.logger.SetFormatter(hook.NewTextFormatter())
		w.logger.SetOutput(colorable.NewColorableStdout())
	}
}

func (w *Write) writeFile(encodeName string, opts ...file.LogOption) {
	info := file.NewFileLog("info.log", opts...).SetLogFile()
	err := file.NewFileLog("error.log", opts...).SetLogFile()
	fm := hook.NewTextFormatter()
	if encodeName == Mode {
		fm = hook.NewJSONFormatter()
	}
	w.logger.AddHook(lfshook.NewHook(
		lfshook.WriterMap{
			logrus.DebugLevel: info,
			logrus.InfoLevel:  info,
			logrus.WarnLevel:  info,
			logrus.ErrorLevel: err,
			logrus.FatalLevel: err,
			logrus.PanicLevel: err,
		},
		fm,
	))
}

func (w *Write) WriteFileAllLog(encodeName string, opts ...file.LogOption) {
	log := file.NewFileLog("log.log", opts...).SetLogFile()
	fm := hook.NewTextFormatter()
	if encodeName == Mode {
		fm = hook.NewJSONFormatter()
	}
	w.logger.AddHook(lfshook.NewHook(
		lfshook.WriterMap{
			logrus.DebugLevel: log,
			logrus.InfoLevel:  log,
			logrus.WarnLevel:  log,
			logrus.ErrorLevel: log,
			logrus.FatalLevel: log,
			logrus.PanicLevel: log,
		},
		fm,
	))
}
