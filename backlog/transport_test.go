package backlog

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewRetryableTransport(t *testing.T) {
	type args struct {
		initialInterval  time.Duration
		maxInterval      time.Duration
		maxRetryAttempts int
		maxJitterMilli   int
	}
	type expected struct {
		value *RetryableTransport
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "basic",
			args: args{
				initialInterval:  1 * time.Second,
				maxInterval:      30 * time.Second,
				maxRetryAttempts: 5,
				maxJitterMilli:   3000,
			},
			expected: expected{
				value: &RetryableTransport{
					Transport:        http.DefaultTransport,
					InitialInterval:  1 * time.Second,
					MaxInterval:      30 * time.Second,
					MaxRetryAttempts: 5,
					MaxJitterMilli:   3000,
				},
			},
		},
		{
			name: "zero",
			args: args{
				initialInterval:  -1,
				maxInterval:      -1,
				maxRetryAttempts: -1,
				maxJitterMilli:   -1,
			},
			expected: expected{
				value: &RetryableTransport{
					Transport:        http.DefaultTransport,
					InitialInterval:  1 * time.Second,
					MaxInterval:      30 * time.Second,
					MaxRetryAttempts: 5,
					MaxJitterMilli:   3000,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := NewRetryableTransport(tt.args.initialInterval, tt.args.maxInterval, tt.args.maxRetryAttempts, tt.args.maxJitterMilli)
			assert.Equal(t, tt.expected.value, actual)
		})
	}
}

func TestRetryableTransport_RoundTrip(t *testing.T) {
	url := "http://example.com"

	type fields struct {
		Transport        http.RoundTripper
		InitialInterval  time.Duration
		MaxInterval      time.Duration
		MaxRetryAttempts int
		MaxJitterMilli   int
	}
	type args struct {
		req *http.Request
	}
	type expected struct {
		isError bool
		status  int
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		expected expected
	}{
		{
			name: "success on first try",
			fields: fields{
				Transport: &mockRoundTripper{
					responses: []*http.Response{
						{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(bytes.NewBufferString("ok")),
							Header:     http.Header{},
						},
					},
				},
				InitialInterval:  1 * time.Millisecond,
				MaxInterval:      10 * time.Millisecond,
				MaxRetryAttempts: 3,
				MaxJitterMilli:   1,
			},
			args: args{
				req: func() *http.Request {
					req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
					return req
				}(),
			},
			expected: expected{
				isError: false,
				status:  http.StatusOK,
			},
		},
		{
			name: "retry then success 500 to 200",
			fields: fields{
				Transport: &mockRoundTripper{
					responses: []*http.Response{
						{
							StatusCode: http.StatusInternalServerError,
							Body:       io.NopCloser(bytes.NewBufferString("err")),
							Header:     http.Header{},
						},
						{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(bytes.NewBufferString("ok")),
							Header:     http.Header{},
						},
					},
				},
				InitialInterval:  1 * time.Millisecond,
				MaxInterval:      10 * time.Millisecond,
				MaxRetryAttempts: 3,
				MaxJitterMilli:   1,
			},
			args: args{
				req: func() *http.Request {
					req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
					return req
				}(),
			},
			expected: expected{
				isError: false,
				status:  http.StatusOK,
			},
		},
		{
			name: "retry always 500 return",
			fields: fields{
				Transport: &mockRoundTripper{
					responses: []*http.Response{
						{
							StatusCode: http.StatusInternalServerError,
							Body:       io.NopCloser(bytes.NewBufferString("err")),
							Header:     http.Header{},
						},
						{
							StatusCode: http.StatusInternalServerError,
							Body:       io.NopCloser(bytes.NewBufferString("err")),
							Header:     http.Header{},
						},
						{
							StatusCode: http.StatusInternalServerError,
							Body:       io.NopCloser(bytes.NewBufferString("err")),
							Header:     http.Header{},
						},
					},
				},
				InitialInterval:  1 * time.Millisecond,
				MaxInterval:      10 * time.Millisecond,
				MaxRetryAttempts: 3,
				MaxJitterMilli:   1,
			},
			args: args{
				req: func() *http.Request {
					req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
					return req
				}(),
			},
			expected: expected{
				isError: true,
				status:  http.StatusInternalServerError,
			},
		},
		{
			name: "429 with ratelimit-reset header",
			fields: fields{
				Transport: &mockRoundTripper{
					responses: []*http.Response{
						{
							StatusCode: http.StatusTooManyRequests,
							Body:       io.NopCloser(bytes.NewBufferString("ratelimit")),
							Header:     http.Header{resetHeaderKey: []string{"1"}},
						},
						{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(bytes.NewBufferString("ok")),
							Header:     http.Header{},
						},
					},
				},
				InitialInterval:  1 * time.Millisecond,
				MaxInterval:      10 * time.Millisecond,
				MaxRetryAttempts: 2,
				MaxJitterMilli:   1,
			},
			args: args{
				req: func() *http.Request {
					req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
					return req
				}(),
			},
			expected: expected{
				isError: false,
				status:  http.StatusOK,
			},
		},
		{
			name: "429 with ratelimit-reset header with zero value",
			fields: fields{
				Transport: &mockRoundTripper{
					responses: []*http.Response{
						{
							StatusCode: http.StatusTooManyRequests,
							Body:       io.NopCloser(bytes.NewBufferString("ratelimit")),
							Header:     http.Header{resetHeaderKey: []string{"0"}},
						},
						{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(bytes.NewBufferString("ok")),
							Header:     http.Header{},
						},
					},
				},
				InitialInterval:  1 * time.Millisecond,
				MaxInterval:      10 * time.Millisecond,
				MaxRetryAttempts: 2,
				MaxJitterMilli:   1,
			},
			args: args{
				req: func() *http.Request {
					req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
					return req
				}(),
			},
			expected: expected{
				isError: false,
				status:  http.StatusOK,
			},
		},
		{
			name: "429 with ratelimit-reset header with negative value",
			fields: fields{
				Transport: &mockRoundTripper{
					responses: []*http.Response{
						{
							StatusCode: http.StatusTooManyRequests,
							Body:       io.NopCloser(bytes.NewBufferString("ratelimit")),
							Header:     http.Header{resetHeaderKey: []string{"-1"}},
						},
						{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(bytes.NewBufferString("ok")),
							Header:     http.Header{},
						},
					},
				},
				InitialInterval:  1 * time.Millisecond,
				MaxInterval:      10 * time.Millisecond,
				MaxRetryAttempts: 2,
				MaxJitterMilli:   1,
			},
			args: args{
				req: func() *http.Request {
					req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
					return req
				}(),
			},
			expected: expected{
				isError: false,
				status:  http.StatusOK,
			},
		},
		{
			name: "429 with ratelimit-reset header with invalid value",
			fields: fields{
				Transport: &mockRoundTripper{
					responses: []*http.Response{
						{
							StatusCode: http.StatusTooManyRequests,
							Body:       io.NopCloser(bytes.NewBufferString("ratelimit")),
							Header:     http.Header{resetHeaderKey: []string{"invalid"}},
						},
					},
				},
				InitialInterval:  1 * time.Millisecond,
				MaxInterval:      10 * time.Millisecond,
				MaxRetryAttempts: 1,
				MaxJitterMilli:   1,
			},
			args: args{
				req: func() *http.Request {
					req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
					return req
				}(),
			},
			expected: expected{
				isError: true,
				status:  http.StatusTooManyRequests,
			},
		},
		{
			name: "transport returns error",
			fields: fields{
				Transport: &mockRoundTripper{
					errors: []error{errors.New("network error")},
				},
				InitialInterval:  1 * time.Millisecond,
				MaxInterval:      10 * time.Millisecond,
				MaxRetryAttempts: 1,
				MaxJitterMilli:   1,
			},
			args: args{
				req: func() *http.Request {
					req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
					return req
				}(),
			},
			expected: expected{
				isError: true,
				status:  0,
			},
		},
		{
			name: "context canceled",
			fields: fields{
				Transport: &mockRoundTripper{
					responses: []*http.Response{
						{
							StatusCode: http.StatusInternalServerError,
							Body:       io.NopCloser(bytes.NewBufferString("err")),
							Header:     http.Header{},
						},
					},
				},
				InitialInterval:  10 * time.Millisecond,
				MaxInterval:      10 * time.Millisecond,
				MaxRetryAttempts: 2,
				MaxJitterMilli:   1,
			},
			args: args{
				req: func() *http.Request {
					ctx, cancel := context.WithCancel(context.Background())
					req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
					go func() {
						time.Sleep(2 * time.Millisecond)
						cancel()
					}()
					return req
				}(),
			},
			expected: expected{
				isError: true,
				status:  http.StatusInternalServerError,
			},
		},
		{
			name: "request body not rewindable",
			fields: fields{
				Transport:        &RetryableTransport{},
				InitialInterval:  1 * time.Millisecond,
				MaxInterval:      10 * time.Millisecond,
				MaxRetryAttempts: 1,
				MaxJitterMilli:   1,
			},
			args: args{
				req: func() *http.Request {
					req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, url, io.NopCloser(bytes.NewBufferString("body")))
					return req
				}(),
			},
			expected: expected{
				isError: true,
				status:  0,
			},
		},
		{
			name: "non-retryable status code",
			fields: fields{
				Transport: &mockRoundTripper{
					responses: []*http.Response{
						{
							StatusCode: http.StatusBadRequest,
							Body:       io.NopCloser(bytes.NewBufferString("bad request")),
							Header:     http.Header{},
						},
					},
				},
				InitialInterval:  1 * time.Millisecond,
				MaxInterval:      10 * time.Millisecond,
				MaxRetryAttempts: 1,
				MaxJitterMilli:   1,
			},
			args: args{
				req: func() *http.Request {
					req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
					return req
				}(),
			},
			expected: expected{
				isError: false,
				status:  http.StatusBadRequest,
			},
		},
		{
			name: "request body is rewindable via GetBody",
			fields: fields{
				Transport: &mockRoundTripper{
					responses: []*http.Response{
						{
							StatusCode: http.StatusInternalServerError,
							Body:       io.NopCloser(bytes.NewBufferString("err")),
							Header:     http.Header{},
						},
						{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(bytes.NewBufferString("ok")),
							Header:     http.Header{},
						},
					},
				},
				InitialInterval:  1 * time.Millisecond,
				MaxInterval:      10 * time.Millisecond,
				MaxRetryAttempts: 2,
				MaxJitterMilli:   1,
			},
			args: args{
				req: func() *http.Request {
					body := []byte("body")
					req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, url, io.NopCloser(bytes.NewBuffer(body)))
					req.GetBody = func() (io.ReadCloser, error) {
						return io.NopCloser(bytes.NewBuffer(body)), nil
					}
					return req
				}(),
			},
			expected: expected{
				isError: false,
				status:  http.StatusOK,
			},
		},
		{
			name: "request body not rewindable with GetBody error",
			fields: fields{
				Transport: &mockRoundTripper{
					responses: []*http.Response{
						{
							StatusCode: http.StatusInternalServerError,
							Body:       io.NopCloser(bytes.NewBufferString("err")),
							Header:     http.Header{},
						},
						{
							StatusCode: http.StatusInternalServerError,
							Body:       io.NopCloser(bytes.NewBufferString("err")),
							Header:     http.Header{},
						},
						{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(bytes.NewBufferString("ok")),
							Header:     http.Header{},
						},
					},
				},
				InitialInterval:  1 * time.Millisecond,
				MaxInterval:      10 * time.Millisecond,
				MaxRetryAttempts: 3,
				MaxJitterMilli:   1,
			},
			args: args{
				req: func() *http.Request {
					req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, url, io.NopCloser(bytes.NewBufferString("body")))
					req.GetBody = func() (io.ReadCloser, error) {
						return nil, errors.New("cannot rewind")
					}
					return req
				}(),
			},
			expected: expected{
				isError: true,
				status:  http.StatusOK,
			},
		},
		{
			name: "response body read error",
			fields: fields{
				Transport: &mockRoundTripper{
					responses: []*http.Response{
						{
							StatusCode: http.StatusInternalServerError,
							Body:       &mockReadCloser{},
							Header:     http.Header{},
						},
					},
				},
				InitialInterval:  1 * time.Millisecond,
				MaxInterval:      10 * time.Millisecond,
				MaxRetryAttempts: 1,
				MaxJitterMilli:   1,
			},
			args: args{
				req: func() *http.Request {
					req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
					return req
				}(),
			},
			expected: expected{
				isError: true,
				status:  http.StatusInternalServerError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &RetryableTransport{
				Transport:        tt.fields.Transport,
				InitialInterval:  tt.fields.InitialInterval,
				MaxInterval:      tt.fields.MaxInterval,
				MaxRetryAttempts: tt.fields.MaxRetryAttempts,
				MaxJitterMilli:   tt.fields.MaxJitterMilli,
			}
			resp, err := o.RoundTrip(tt.args.req)
			if tt.expected.isError {
				assert.Error(t, err)
				if resp != nil && tt.expected.status != 0 {
					assert.Equal(t, tt.expected.status, resp.StatusCode)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				if tt.expected.status != 0 && resp != nil {
					assert.Equal(t, tt.expected.status, resp.StatusCode)
				}
			}
		})
	}
}
