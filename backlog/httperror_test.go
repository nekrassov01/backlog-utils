// backlog/backlog_test.go
package backlog

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetErrorMessage(t *testing.T) {
	type args struct {
		resp *http.Response
	}
	type expected struct {
		value string
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "basic",
			args: args{
				resp: &http.Response{
					StatusCode: http.StatusBadRequest,
					Body:       io.NopCloser(bytes.NewBufferString(`{"errors": [{"message": "Invalid request", "code": 400}]}`)),
				},
			},
			expected: expected{
				value: "Invalid request",
			},
		},
		{
			name: "multiple error messages",
			args: args{
				resp: &http.Response{
					StatusCode: http.StatusBadRequest,
					Body:       io.NopCloser(bytes.NewBufferString(`{"errors": [{"message": "Invalid request", "code": 400},{"message": "Missing parameter", "code": 401}]}`)),
				},
			},
			expected: expected{
				value: "Invalid request; Missing parameter",
			},
		},
		{
			name: "empty",
			args: args{
				resp: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBufferString("")),
				},
			},
			expected: expected{
				value: http.StatusText(http.StatusInternalServerError),
			},
		},
		{
			name: "invalid json",
			args: args{
				resp: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBufferString("{invalid json")),
				},
			},
			expected: expected{
				value: http.StatusText(http.StatusInternalServerError),
			},
		},
		{
			name: "no errors field",
			args: args{
				resp: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBufferString(`{"message": "error"}`)),
				},
			},
			expected: expected{
				value: http.StatusText(http.StatusInternalServerError),
			},
		},
		{
			name: "nil",
			args: args{
				resp: nil,
			},
			expected: expected{
				value: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := GetErrorMessage(tt.args.resp)
			assert.Equal(t, tt.expected.value, actual)
		})
	}
}
