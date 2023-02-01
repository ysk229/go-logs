package logrus

import (
	"github.com/mattn/go-colorable"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"

	log "github.com/ysk229/go-logs"
	"github.com/ysk229/go-logs/contrib/logrus/hook"
	"github.com/ysk229/go-logs/file"
)

type Write struct {
	logger *logrus.Logger
}

func NewWrite(logger *logrus.Logger) *Write {
	return &Write{logger: logger}
}

func (w *Write) setConsoleFormatter(encodeName string) {
	if encodeName == "json" {
		w.logger.SetFormatter(hook.NewJsonFormatter())
		w.logger.SetOutput(log.NewJSONColorable())
	} else {
		w.logger.SetFormatter(hook.NewTextFormatter())
		w.logger.SetOutput(colorable.NewColorableStdout())
	}
}

func (w *Write) writeFile(encodeName string, opts ...file.FileLogOption) {
	info := file.NewFileLog("info.log", opts...).SetLogFile()
	err := file.NewFileLog("error.log", opts...).SetLogFile()
	fm := hook.NewTextFormatter()
	if encodeName == "json" {
		fm = hook.NewJsonFormatter()
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

func (w *Write) writeFileAllLog(encodeName string, opts ...file.FileLogOption) {
	log := file.NewFileLog("log.log", opts...).SetLogFile()
	fm := hook.NewTextFormatter()
	if encodeName == "json" {
		fm = hook.NewJsonFormatter()
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
