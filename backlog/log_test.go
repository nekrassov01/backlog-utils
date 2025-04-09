package backlog

import (
	"bytes"
	"log/slog"
	"reflect"
	"testing"
)

func TestNewLogger(t *testing.T) {
	type args struct {
		level string
	}
	tests := []struct {
		name string
		args args
		want *slog.Logger
	}{
		{
			name: "lowercase",
			args: args{
				level: "info",
			},
			want: slog.New(slog.NewTextHandler(&bytes.Buffer{},
				&slog.HandlerOptions{
					Level: slog.LevelInfo,
				}),
			),
		},
		{
			name: "uppercase",
			args: args{
				level: "INFO",
			},
			want: slog.New(slog.NewTextHandler(&bytes.Buffer{},
				&slog.HandlerOptions{
					Level: slog.LevelInfo,
				}),
			),
		},
		{
			name: "debug",
			args: args{
				level: "debug",
			},
			want: slog.New(slog.NewTextHandler(&bytes.Buffer{},
				&slog.HandlerOptions{
					Level: slog.LevelDebug,
				}),
			),
		},
		{
			name: "warn",
			args: args{
				level: "warn",
			},
			want: slog.New(slog.NewTextHandler(&bytes.Buffer{},
				&slog.HandlerOptions{
					Level: slog.LevelWarn,
				}),
			),
		},
		{
			name: "error",
			args: args{
				level: "error",
			},
			want: slog.New(slog.NewTextHandler(&bytes.Buffer{},
				&slog.HandlerOptions{
					Level: slog.LevelError,
				}),
			),
		},
		{
			name: "fallback",
			args: args{
				level: "dummy",
			},
			want: slog.New(slog.NewTextHandler(&bytes.Buffer{},
				&slog.HandlerOptions{
					Level: slog.LevelInfo,
				}),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if got := NewLogger(w, tt.args.level); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}
