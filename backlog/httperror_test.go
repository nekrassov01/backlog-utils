// backlog/backlog_test.go
package backlog

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

func TestGetErrorMessage(t *testing.T) {
	type args struct {
		resp *http.Response
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "basic",
			args: args{
				resp: &http.Response{
					StatusCode: http.StatusBadRequest,
					Body:       io.NopCloser(bytes.NewBufferString(`{"errors": [{"message": "Invalid request", "code": 400}]}`)),
				},
			},
			want: "Invalid request",
		},
		{
			name: "multiple error messages",
			args: args{
				resp: &http.Response{
					StatusCode: http.StatusBadRequest,
					Body:       io.NopCloser(bytes.NewBufferString(`{"errors": [{"message": "Invalid request", "code": 400},{"message": "Missing parameter", "code": 401}]}`)),
				},
			},
			want: "Invalid request; Missing parameter",
		},
		{
			name: "empty",
			args: args{
				resp: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBufferString("")),
				},
			},
			want: "",
		},
		{
			name: "invalid json",
			args: args{
				resp: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBufferString("{invalid json")),
				},
			},
			want: "",
		},
		{
			name: "no errors field",
			args: args{
				resp: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBufferString(`{"message": "error"}`)),
				},
			},
			want: "",
		},
		{
			name: "nil",
			args: args{
				resp: nil,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetErrorMessage(tt.args.resp); got != tt.want {
				t.Errorf("GetErrorMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
