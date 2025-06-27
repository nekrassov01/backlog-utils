package log

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	type args struct {
		level string
	}
	type expected struct {
		value *slog.Logger
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "lowercase",
			args: args{
				level: "info",
			},
			expected: expected{
				value: slog.New(slog.NewTextHandler(&bytes.Buffer{},
					&slog.HandlerOptions{
						Level: slog.LevelInfo,
					}),
				),
			},
		},
		{
			name: "uppercase",
			args: args{
				level: "INFO",
			},
			expected: expected{
				value: slog.New(slog.NewTextHandler(&bytes.Buffer{},
					&slog.HandlerOptions{
						Level: slog.LevelInfo,
					}),
				),
			},
		},
		{
			name: "debug",
			args: args{
				level: "debug",
			},
			expected: expected{
				value: slog.New(slog.NewTextHandler(&bytes.Buffer{},
					&slog.HandlerOptions{
						Level: slog.LevelDebug,
					}),
				),
			},
		},
		{
			name: "warn",
			args: args{
				level: "warn",
			},
			expected: expected{
				value: slog.New(slog.NewTextHandler(&bytes.Buffer{},
					&slog.HandlerOptions{
						Level: slog.LevelWarn,
					}),
				),
			},
		},
		{
			name: "error",
			args: args{
				level: "error",
			},
			expected: expected{
				value: slog.New(slog.NewTextHandler(&bytes.Buffer{},
					&slog.HandlerOptions{
						Level: slog.LevelError,
					}),
				),
			},
		},
		{
			name: "fallback",
			args: args{
				level: "dummy",
			},
			expected: expected{
				value: slog.New(slog.NewTextHandler(&bytes.Buffer{},
					&slog.HandlerOptions{
						Level: slog.LevelInfo,
					}),
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			actual := NewLogger(w, tt.args.level)
			assert.Equal(t, tt.expected.value, actual)
		})
	}
}
