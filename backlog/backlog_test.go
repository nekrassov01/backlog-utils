// backlog/backlog_test.go
package backlog

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

func TestBacklog_Do(t *testing.T) {
	type fields struct {
		Backlog *Backlog
	}
	type args struct {
		req *http.Request
	}
	type mock struct {
		status       int
		body         string
		retryHeaders map[string]string
		retryCount   int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    mock
		want    *http.Response
		wantErr bool
	}{
		{
			name: "basic",
			fields: fields{
				Backlog: &Backlog{
					Writer:           &bytes.Buffer{},
					BaseURL:          "http://example.com",
					APIKey:           "dummy",
					MaxRetryAttempts: 2,
					MaxJitterMilli:   10,
				},
			},
			args: args{
				req: mustNewRequest(http.MethodGet, "http://example.com/test", nil),
			},
			mock: mock{
				status: http.StatusOK,
				body:   http.StatusText(http.StatusOK),
			},
			want: &http.Response{
				StatusCode: http.StatusOK,
			},
			wantErr: false,
		},
		{
			name: "retry",
			fields: fields{
				Backlog: &Backlog{
					Writer:           &bytes.Buffer{},
					BaseURL:          "http://example.com",
					APIKey:           "dummy",
					MaxRetryAttempts: 2,
					MaxJitterMilli:   10,
				},
			},
			args: args{
				req: mustNewRequest(http.MethodGet, "http://example.com/test", nil),
			},
			mock: mock{
				status: http.StatusOK,
				body:   http.StatusText(http.StatusOK),
				retryHeaders: map[string]string{
					"X-RateLimit-Reset": strconv.FormatInt(time.Now().Add(1*time.Second).Unix(), 10),
				},
				retryCount: 1,
			},
			want: &http.Response{
				StatusCode: http.StatusOK,
			},
			wantErr: false,
		},
		{
			name: "internal server error",
			fields: fields{
				Backlog: &Backlog{
					Writer:           &bytes.Buffer{},
					BaseURL:          "http://example.com",
					APIKey:           "dummy",
					MaxRetryAttempts: 2,
					MaxJitterMilli:   10,
				},
			},
			args: args{
				req: mustNewRequest(http.MethodGet, "http://example.com/test", nil),
			},
			mock: mock{
				status: http.StatusInternalServerError,
				body:   http.StatusText(http.StatusInternalServerError),
			},
			want: &http.Response{
				StatusCode: http.StatusInternalServerError,
			},
			wantErr: false,
		},
		{
			name: "exceeded max retries",
			fields: fields{
				Backlog: &Backlog{
					Writer:           &bytes.Buffer{},
					BaseURL:          "http://example.com",
					APIKey:           "dummy",
					MaxRetryAttempts: 2,
					MaxJitterMilli:   10,
				},
			},
			args: args{
				req: mustNewRequest(http.MethodGet, "http://example.com/test", nil),
			},
			mock: mock{
				status: http.StatusTooManyRequests,
				body:   http.StatusText(http.StatusTooManyRequests),
				retryHeaders: map[string]string{
					"X-RateLimit-Reset": strconv.FormatInt(time.Now().Add(1*time.Second).Unix(), 10),
				},
				retryCount: 3, // Exceeding the max retry attempts
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := tt.fields.Backlog
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			n := 0
			httpmock.RegisterResponder(
				http.MethodGet,
				fmt.Sprintf("%s/test", o.BaseURL),
				func(req *http.Request) (*http.Response, error) {
					if n < tt.mock.retryCount {
						n++
						resp := httpmock.NewStringResponse(http.StatusTooManyRequests, tt.mock.body)
						for k, v := range tt.mock.retryHeaders {
							resp.Header.Set(k, v)
						}
						return resp, nil
					}
					return httpmock.NewStringResponder(tt.mock.status, tt.mock.body)(req)
				},
			)
			got, err := o.Do(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Backlog.Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want == nil {
				if got != nil {
					t.Errorf("Backlog.Do() got = %v, want nil", got)
				}
				return
			}
			if got.StatusCode != tt.want.StatusCode {
				t.Errorf("Backlog.Do() status = %v, want %v", got.StatusCode, tt.want.StatusCode)
			}
		})
	}
}
