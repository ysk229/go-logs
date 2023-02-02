package zap

import (
	"testing"

	log "github.com/ysk229/go-logs"
	"github.com/ysk229/go-logs/config"
)

func Test_zapLogger_Log(t *testing.T) {
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
