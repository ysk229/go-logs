package zap

import (
	log "github.com/ysk229/go-logs"
	"github.com/ysk229/go-logs/contrib/zap/encoder"
	"github.com/ysk229/go-logs/file"

	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const Mode = "json"

const bufferSize = 4096

type Write struct {
	confLevel string
}

func NewWrite(level string) *Write {
	return &Write{confLevel: level}
}

func Level(confLevel string) zapcore.Level {
	// 设置级别
	logLevel := zap.WarnLevel
	switch confLevel {
	case "debug":
		logLevel = zap.DebugLevel
	case "info":
		logLevel = zap.InfoLevel
	case "warn":
		logLevel = zap.WarnLevel
	case "error":
		logLevel = zap.ErrorLevel
	case "panic":
		logLevel = zap.PanicLevel
	case "fatal":
		logLevel = zap.FatalLevel
	}
	return logLevel
}

func (w *Write) SetLevel(level string) {
	w.confLevel = level
}

func (w *Write) setFileEncodeName(encodeName string) zapcore.Encoder {
	if encodeName == Mode {
		return encoder.NewJSONEncoder()
	}
	return encoder.NewTextNoColorEncoder()
}

func (w *Write) setConsoleEncodeName(encodeName string) zapcore.Encoder {
	if encodeName == Mode {
		return encoder.NewJSONEncoder()
	}
	return encoder.NewTextColorEncoder()
}

func (w *Write) SetFileInfo() zap.LevelEnablerFunc {
	return zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev < zap.ErrorLevel && lev >= zap.DebugLevel && lev >= Level(w.confLevel)
	})
}

func (w *Write) SetFileError() zap.LevelEnablerFunc {
	return zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev >= zap.ErrorLevel && lev >= Level(w.confLevel)
	})
}

func (w *Write) SetFileAll() zap.LevelEnablerFunc {
	return zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev >= Level(w.confLevel)
	})
}

func (w *Write) writeFile(encodeName string, lev zapcore.LevelEnabler, fileName string,
	opts ...file.LogOption,
) zapcore.Core {
	return zapcore.NewCore(
		w.setFileEncodeName(encodeName),
		zapcore.AddSync(file.NewFileLog(fileName, opts...).SetLogFile()),
		lev,
	)
}

// WriteAsyncFile 异步落盘
func (w *Write) WriteAsyncFile(encodeName string, lev zapcore.LevelEnabler, fileName string,
	opts ...file.LogOption,
) zapcore.Core {
	return zapcore.NewCore(
		w.setFileEncodeName(encodeName),
		&zapcore.BufferedWriteSyncer{
			WS:   zapcore.AddSync(file.NewFileLog(fileName, opts...).SetLogFile()),
			Size: bufferSize,
		},
		lev,
	)
}

func (w *Write) writeConsole(encodeName string, lev zapcore.LevelEnabler) zapcore.Core {
	io := colorable.NewColorableStdout()
	if encodeName == Mode {
		io = log.NewJSONColorable()
	}
	return zapcore.NewCore(
		w.setConsoleEncodeName(encodeName),
		zapcore.AddSync(io),
		lev,
	)
}
