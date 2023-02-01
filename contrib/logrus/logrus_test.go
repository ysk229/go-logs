package logrus

import (
	"bytes"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	log "github.com/ysk229/go-logs"
	"github.com/ysk229/go-logs/config"
)

func TestLoggerLog(t *testing.T) {
	tests := map[string]struct {
		level     logrus.Level
		formatter logrus.Formatter
		logLevel  log.Level
		kvs       []interface{}
		want      string
	}{
		"json format": {
			level:     logrus.InfoLevel,
			formatter: &logrus.JSONFormatter{},
			logLevel:  log.LevelInfo,
			kvs:       []interface{}{"case", "json format", "msg", "1"},
			want:      `{"case":"json format","level":"info","msg":"1"`,
		},
		"level unmatch": {
			level:     logrus.InfoLevel,
			formatter: &logrus.JSONFormatter{},
			logLevel:  log.LevelDebug,
			kvs:       []interface{}{"case", "level unmatch", "msg", "1"},
			want:      "",
		},
		"no tags": {
			level:     logrus.InfoLevel,
			formatter: &logrus.JSONFormatter{},
			logLevel:  log.LevelInfo,
			kvs:       []interface{}{"msg", "1"},
			want:      `{"level":"info","msg":"1"`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			output := new(bytes.Buffer)
			logger := logrus.New()
			logger.Level = test.level
			logger.Out = output
			logger.Formatter = test.formatter
			wrapped := &Logger{logrus: logger}
			_ = wrapped.Log(test.logLevel, test.kvs...)

			assert.True(t, strings.HasPrefix(output.String(), test.want))
		})
	}
}

func TestLogger_Log(t *testing.T) {
	type args struct {
		level   log.Level
		keyvals []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "info",
			args: args{
				level:   0,
				keyvals: []interface{}{"info-key", "info-value"},
			},
		},
		{
			name: "warn",
			args: args{
				level:   1,
				keyvals: []interface{}{"warn-key", "warn-value", "msg", "test"},
			},
		},
	}

	l := New(&config.Config{Level: "info"})
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if err := l.Log(tt.args.level, tt.args.keyvals...); (err != nil) != tt.wantErr {
					t.Errorf("Log() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}
