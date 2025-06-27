// backlog/backlog_test.go
package backlog

import (
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithWriter(t *testing.T) {
	type args struct {
		w io.Writer
	}
	type expected struct {
		value *Client
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "basic",
			args: args{
				w: io.Discard,
			},
			expected: expected{
				value: &Client{
					Writer: io.Discard,
				},
			},
		},
		{
			name: "nil",
			args: args{
				w: nil,
			},
			expected: expected{
				value: &Client{
					Writer: os.Stdout,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{HTTPClient: &http.Client{}}
			WithWriter(tt.args.w)(client)
			assert.Equal(t, tt.expected.value.Writer, client.Writer)
		})
	}
}

func TestWithTransport(t *testing.T) {
	type args struct {
		transport http.RoundTripper
	}
	type expected struct {
		value *Client
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "basic",
			args: args{
				transport: NewRetryableTransport(1*time.Second, 30*time.Second, 5, 3000),
			},
			expected: expected{
				value: &Client{
					HTTPClient: &http.Client{
						Transport: NewRetryableTransport(1*time.Second, 30*time.Second, 5, 3000),
					},
				},
			},
		},
		{
			name: "nil",
			args: args{
				transport: nil,
			},
			expected: expected{
				value: &Client{
					HTTPClient: &http.Client{
						Transport: http.DefaultTransport,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{HTTPClient: &http.Client{}}
			WithTransport(tt.args.transport)(client)
			assert.Equal(t, tt.expected.value.HTTPClient.Transport, client.HTTPClient.Transport)
		})
	}
}

func TestNewClient(t *testing.T) {
	type args struct {
		url    string
		apiKey string
		opts   []ClientOption
	}
	type expected struct {
		value   *Client
		isError bool
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "basic",
			args: args{
				url:    "https://example.com",
				apiKey: "dummy",
				opts: []ClientOption{
					WithWriter(io.Discard),
					WithTransport(http.DefaultTransport),
				},
			},
			expected: expected{
				value: &Client{
					Writer:  io.Discard,
					BaseURL: "https://example.com",
					APIKey:  "dummy",
					HTTPClient: &http.Client{
						Transport: http.DefaultTransport,
					},
				},
				isError: false,
			},
		},
		{
			name: "empty url",
			args: args{
				url:    "",
				apiKey: "dummy",
				opts:   nil,
			},
			expected: expected{
				value:   nil,
				isError: true,
			},
		},
		{
			name: "empty api key",
			args: args{
				url:    "https://example.com",
				apiKey: "",
				opts:   nil,
			},
			expected: expected{
				value:   nil,
				isError: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := NewClient(tt.args.url, tt.args.apiKey, tt.args.opts...)
			if tt.expected.isError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.value, actual)
		})
	}
}
