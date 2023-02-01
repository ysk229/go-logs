package file

import (
	"io"
	"time"

	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/natefinch/lumberjack"
)

// FileLogOption is fileLog option.
type FileLogOption func(*FileLog)

type FileLog struct {
	file   string
	mode   string // date ,size
	maxAge int    // 90 日志文件存储最大天数
	size   int    // 30,    //日志的最大大小（M）
	path   string // default /var/log/
}

func NewFileLog(name string, opts ...FileLogOption) *FileLog {
	options := FileLog{file: name, mode: "size", maxAge: 90, size: 30, path: "./logs"}

	for _, o := range opts {
		o(&options)
	}
	// options.file = options.path + "/" + options.file //指定日志存储位置
	return &options
}

// Mode file mode
func Mode(mode string) FileLogOption {
	return func(opts *FileLog) {
		if len(mode) > 0 {
			opts.mode = mode
		}
	}
}

// Path path
func Path(path string) FileLogOption {
	return func(opts *FileLog) {
		if len(path) > 0 {
			opts.path = path
		}
	}
}

// MaxAge maxAge
func MaxAge(maxAge int) FileLogOption {
	return func(opts *FileLog) {
		if maxAge > 0 {
			opts.maxAge = maxAge
		}
	}
}

// Size file size
func Size(size int) FileLogOption {
	return func(opts *FileLog) {
		if size > 0 {
			opts.size = size
		}
	}
}

func (f *FileLog) setLogFileNameLumberjack() io.Writer {
	return &lumberjack.Logger{
		Filename:   f.path + "/" + f.file, // 指定日志存储位置
		MaxSize:    f.size,                // 30,    //日志的最大大小（M）
		MaxBackups: 2,                     // 日志的最大保存数量
		MaxAge:     f.maxAge,              // 日志文件存储最大天数
		Compress:   false,                 // 是否执行压缩
	}
}

func (f *FileLog) setLogFileNameRotate() io.Writer {
	writer, _ := rotatelogs.New(f.path+"/%Y%m/%d/"+f.file,
		// rotatelogs.WithLinkName(f.file),                           // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(time.Duration(f.maxAge*24)*time.Hour), // 文件最大保存时间
		rotatelogs.WithRotationTime(24*time.Hour),                   // 日志切割时间间隔
		rotatelogs.WithRotationSize(int64(f.size)*1024*1024),        // 设置文件大小,当大于这个容量时，创建新的日志文件
	)
	return writer
}

func (f *FileLog) SetLogFile() io.Writer {
	if f.mode == "date" {
		return f.setLogFileNameRotate()
	}
	return f.setLogFileNameLumberjack()
}
