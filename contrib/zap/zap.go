package zap

import (
	"fmt"
	log2 "log"

	log "github.com/ysk229/go-logs"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/ysk229/go-logs/config"
	"github.com/ysk229/go-logs/file"
)

var _ log.Logger = (*Logger)(nil)

type Logger struct {
	zap *zap.Logger
	w   *Write
}

func New(conf *config.Config) *Logger {
	w := NewWrite(conf.Level)
	l := getLog(conf, w)
	zap.RedirectStdLog(l)
	zap.ReplaceGlobals(l)
	log2.SetFlags(log2.Lshortfile)

	return &Logger{zap: l, w: w}
}

func (l *Logger) Log(level log.Level, keyvals ...interface{}) error {
	if len(keyvals) == 0 || len(keyvals)%2 != 0 {
		l.zap.Warn(fmt.Sprint("Keyvalues must appear in pairs : ", keyvals))
		return nil
	}
	if len(keyvals) == 0 {
		return nil
	}
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, "")
	}
	var data []zap.Field
	msg := ""
	for i := 0; i < len(keyvals); i += 2 {
		key, ok := keyvals[i].(string)
		if !ok {
			continue
		}
		if key == "msg" {
			msg, _ = keyvals[i+1].(string)
			continue
		}
		data = append(data, zap.Any(fmt.Sprint(keyvals[i]), keyvals[i+1]))
	}
	switch level {
	case log.LevelDebug:
		l.zap.Debug(msg, data...)
	case log.LevelInfo:
		l.zap.Info(msg, data...)
	case log.LevelWarn:
		l.zap.Warn(msg, data...)
	case log.LevelError:
		l.zap.Error(msg, data...)
	case log.LevelFatal:
		l.zap.Fatal(msg, data...)
	}
	return nil
}

func (l *Logger) Sync() error {
	return l.zap.Sync()
}

func (l *Logger) SetLevel(level string) {
	l.w.SetLevel(level)
}

func (l *Logger) Close() error {
	return l.Sync()
}

func getLog(conf *config.Config, w *Write) *zap.Logger {
	lev := "warn"

	both := ""
	format := "text"
	var fileOpts []file.LogOption
	if conf != nil {
		lev = conf.Level
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
		both = conf.Both
		format = conf.Format
	}
	w.SetLevel(lev)
	switch both {
	case "all":
		return zap.New(
			zapcore.NewTee(
				w.writeFile(format, w.SetFileInfo(), "info.log", fileOpts...),
				w.writeFile(format, w.SetFileError(), "error.log", fileOpts...),
				w.writeConsole(format, w.SetFileAll()),
			),
			zap.AddStacktrace(
				zap.NewAtomicLevelAt(zapcore.ErrorLevel),
			),
		)
	case "file":
		return zap.New(
			zapcore.NewTee(
				w.writeFile(format, w.SetFileInfo(), "info.log", fileOpts...),
				w.writeFile(format, w.SetFileError(), "error.log", fileOpts...),
			),
			zap.AddStacktrace(
				zap.NewAtomicLevelAt(zapcore.ErrorLevel),
			),
		)
	}
	return zap.New(
		zapcore.NewTee(
			w.writeConsole(format, w.SetFileAll()),
		),
		zap.AddStacktrace(
			zap.NewAtomicLevelAt(zapcore.ErrorLevel),
		),
	)
}
