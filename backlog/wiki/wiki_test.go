package wiki

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/nekrassov01/backlog-utils/backlog"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	type args struct {
		url    string
		apiKey string
		opts   []backlog.ClientOption
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
				opts: []backlog.ClientOption{
					backlog.WithWriter(io.Discard),
					backlog.WithTransport(http.DefaultTransport),
				},
			},
			expected: expected{
				value: &Client{
					&backlog.Client{
						Writer:  io.Discard,
						BaseURL: "https://example.com",
						APIKey:  "dummy",
						HTTPClient: &http.Client{
							Transport: http.DefaultTransport,
						},
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

func TestWiki_List(t *testing.T) {
	type fields struct {
		Backlog *backlog.Client
	}
	type args struct {
		projectKey string
		pattern    string
	}
	type expected struct {
		value   []*Page
		isError bool
	}
	type mock struct {
		status int
		body   string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		mock     mock
		expected expected
	}{
		{
			name: "basic",
			fields: fields{
				Backlog: &backlog.Client{
					Writer:     io.Discard,
					BaseURL:    "https://example.com",
					APIKey:     "dummy",
					HTTPClient: &http.Client{},
				},
			},
			args: args{
				projectKey: "dummy",
				pattern:    "",
			},
			expected: expected{
				value: []*Page{
					{
						ID:        1,
						ProjectID: 123,
						Name:      "Test Page",
						Content:   "",
					},
				},
				isError: false,
			},
			mock: mock{
				status: 200,
				body:   `[{"id":1,"projectId":123,"name":"Test Page","content":""}]`,
			},
		},
		{
			name: "pattern",
			fields: fields{
				Backlog: &backlog.Client{
					Writer:     io.Discard,
					BaseURL:    "https://example.com",
					APIKey:     "dummy",
					HTTPClient: &http.Client{},
				},
			},
			args: args{
				projectKey: "dummy",
				pattern:    "Test",
			},
			expected: expected{
				value: []*Page{
					{
						ID:        1,
						ProjectID: 123,
						Name:      "Test Page",
					},
				},
				isError: false,
			},
			mock: mock{
				status: 200,
				body:   `[{"id":1,"projectId":123,"name":"Test Page"},{"id":2,"projectId":123,"name":"Other Page"}]`,
			},
		},
		{
			name: "invalid pattern",
			fields: fields{
				Backlog: &backlog.Client{
					Writer:     io.Discard,
					BaseURL:    "https://example.com",
					APIKey:     "dummy",
					HTTPClient: &http.Client{},
				},
			},
			args: args{
				projectKey: "dummy",
				pattern:    "[",
			},
			expected: expected{
				value:   nil,
				isError: true,
			},
			mock: mock{
				status: 0,
				body:   "",
			},
		},
		{
			name: "empty url",
			fields: fields{
				Backlog: &backlog.Client{
					Writer:     io.Discard,
					BaseURL:    "",
					APIKey:     "dummy",
					HTTPClient: &http.Client{},
				},
			},
			args: args{
				projectKey: "dummy",
				pattern:    "",
			},
			expected: expected{
				value:   nil,
				isError: true,
			},
			mock: mock{
				status: 0,
				body:   "",
			},
		},
		{
			name: "empty api key",
			fields: fields{
				Backlog: &backlog.Client{
					Writer:     io.Discard,
					BaseURL:    "https://example.com",
					APIKey:     "",
					HTTPClient: &http.Client{},
				},
			},
			args: args{
				projectKey: "dummy",
				pattern:    "",
			},
			expected: expected{
				value:   nil,
				isError: true,
			},
			mock: mock{
				status: 0,
				body:   "",
			},
		},
		{
			name: "empty project key",
			fields: fields{
				Backlog: &backlog.Client{
					Writer:     io.Discard,
					BaseURL:    "https://example.com",
					APIKey:     "dummy",
					HTTPClient: &http.Client{},
				},
			},
			args: args{
				projectKey: "",
				pattern:    "",
			},
			expected: expected{
				value:   nil,
				isError: true,
			},
			mock: mock{
				status: 0,
				body:   "",
			},
		},
		{
			name: "api error",
			fields: fields{
				Backlog: &backlog.Client{
					Writer:     io.Discard,
					BaseURL:    "https://example.com",
					APIKey:     "dummy",
					HTTPClient: &http.Client{},
				},
			},
			args: args{
				projectKey: "dummy",
				pattern:    "",
			},
			expected: expected{
				value:   nil,
				isError: true,
			},
			mock: mock{
				status: 500,
				body:   `{"errors":[{"message":"Internal Server Error"}]}`,
			},
		},
		{
			name: "invalid response",
			fields: fields{
				Backlog: &backlog.Client{
					Writer:     io.Discard,
					BaseURL:    "https://example.com",
					APIKey:     "dummy",
					HTTPClient: &http.Client{},
				},
			},
			args: args{
				projectKey: "dummy",
				pattern:    "",
			},
			expected: expected{
				value:   nil,
				isError: true,
			},
			mock: mock{
				status: 200,
				body:   `[{"id":1,"projectId":123,"name":}]`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Client{
				Client: tt.fields.Backlog,
			}
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			if tt.mock.status != 0 {
				httpmock.RegisterResponder(
					http.MethodGet,
					fmt.Sprintf("%s/api/v2/wikis?projectIdOrKey=%s&apiKey=%s", o.BaseURL, tt.args.projectKey, o.APIKey),
					httpmock.NewStringResponder(tt.mock.status, tt.mock.body),
				)
			}
			actual, err := o.List(tt.args.projectKey, tt.args.pattern)
			if tt.expected.isError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.value, actual)
		})
	}
}

func TestWiki_Get(t *testing.T) {
	type fields struct {
		Backlog *backlog.Client
	}
	type args struct {
		id int64
	}
	type expected struct {
		value   *Page
		isError bool
	}
	type mock struct {
		status int
		body   string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		expected expected
		mock     mock
	}{
		{
			name: "basic",
			fields: fields{
				Backlog: &backlog.Client{
					Writer:     io.Discard,
					BaseURL:    "https://example.com",
					APIKey:     "dummy",
					HTTPClient: &http.Client{},
				},
			},
			args: args{
				id: 1,
			},
			expected: expected{
				value: &Page{
					ID:        1,
					ProjectID: 123,
					Name:      "Test Page",
					Content:   "Sample Content",
				},
				isError: false,
			},
			mock: mock{
				status: 200,
				body:   `{"id":1,"projectId":123,"name":"Test Page","content":"Sample Content"}`,
			},
		},
		{
			name: "empty url",
			fields: fields{
				Backlog: &backlog.Client{
					Writer:     io.Discard,
					BaseURL:    "",
					APIKey:     "dummy",
					HTTPClient: &http.Client{},
				},
			},
			args: args{
				id: 1,
			},
			expected: expected{
				value:   nil,
				isError: true,
			},
			mock: mock{
				status: 0,
				body:   "",
			},
		},
		{
			name: "empty api key",
			fields: fields{
				Backlog: &backlog.Client{
					Writer:     io.Discard,
					BaseURL:    "https://example.com",
					APIKey:     "",
					HTTPClient: &http.Client{},
				},
			},
			args: args{
				id: 1,
			},
			expected: expected{
				value:   nil,
				isError: true,
			},
			mock: mock{
				status: 0,
				body:   "",
			},
		},
		{
			name: "invalid id",
			fields: fields{
				Backlog: &backlog.Client{
					Writer:     io.Discard,
					BaseURL:    "https://example.com",
					APIKey:     "dummy",
					HTTPClient: &http.Client{},
				},
			},
			args: args{
				id: -1,
			},
			expected: expected{
				value:   nil,
				isError: true,
			},
			mock: mock{
				status: 0,
				body:   "",
			},
		},
		{
			name: "api error",
			fields: fields{
				Backlog: &backlog.Client{
					Writer:     io.Discard,
					BaseURL:    "https://example.com",
					APIKey:     "dummy",
					HTTPClient: &http.Client{},
				},
			},
			args: args{
				id: 1,
			},
			expected: expected{
				value:   nil,
				isError: true,
			},
			mock: mock{
				status: 404,
				body:   `{"errors":[{"message":"Not Found"}]}`,
			},
		},
		{
			name: "invalid response",
			fields: fields{
				Backlog: &backlog.Client{
					Writer:     io.Discard,
					BaseURL:    "https://example.com",
					APIKey:     "dummy",
					HTTPClient: &http.Client{},
				},
			},
			args: args{
				id: 1,
			},
			expected: expected{
				value:   nil,
				isError: true,
			},
			mock: mock{
				status: 200,
				body:   `{"id":1,"projectId":123,"name":}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Client{
				Client: tt.fields.Backlog,
			}
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			if tt.mock.status != 0 {
				httpmock.RegisterResponder(
					http.MethodGet,
					fmt.Sprintf("%s/api/v2/wikis/%d?apiKey=%s", o.BaseURL, tt.args.id, o.APIKey),
					httpmock.NewStringResponder(tt.mock.status, tt.mock.body),
				)
			}
			actual, err := o.Get(tt.args.id)
			if tt.expected.isError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.value, actual)
		})
	}
}

func TestWiki_Rename(t *testing.T) {
	type fields struct {
		Backlog *backlog.Client
	}
	type args struct {
		page *Page
		old  string
		new  string
	}
	type expected struct {
		isError bool
	}
	type mock struct {
		status int
		body   string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		expected expected
		mock     mock
	}{
		{
			name: "basic",
			fields: fields{
				Backlog: &backlog.Client{
					Writer:     io.Discard,
					BaseURL:    "https://example.com",
					APIKey:     "dummy",
					HTTPClient: &http.Client{},
				},
			},
			args: args{
				page: &Page{
					ID:   1,
					Name: "Old Name",
				},
				old: "Old",
				new: "New",
			},
			expected: expected{
				isError: false,
			},
			mock: mock{
				status: 200,
				body:   "",
			},
		},
		{
			name: "empty page",
			fields: fields{
				Backlog: &backlog.Client{
					Writer:     io.Discard,
					BaseURL:    "https://example.com",
					APIKey:     "dummy",
					HTTPClient: &http.Client{},
				},
			},
			args: args{
				page: nil,
				old:  "Old",
				new:  "New",
			},
			expected: expected{
				isError: true,
			},
			mock: mock{
				status: 0,
				body:   "",
			},
		},
		{
			name: "empty old string",
			fields: fields{
				Backlog: &backlog.Client{
					Writer:     io.Discard,
					BaseURL:    "https://example.com",
					APIKey:     "dummy",
					HTTPClient: &http.Client{},
				},
			},
			args: args{
				page: &Page{
					ID:   1,
					Name: "Old Name",
				},
				old: "",
				new: "new",
			},
			expected: expected{
				isError: true,
			},
			mock: mock{
				status: 0,
				body:   "",
			},
		},
		{
			name: "empty new string",
			fields: fields{
				Backlog: &backlog.Client{
					Writer:     io.Discard,
					BaseURL:    "https://example.com",
					APIKey:     "dummy",
					HTTPClient: &http.Client{},
				},
			},
			args: args{
				page: &Page{
					ID:   1,
					Name: "Old Name",
				},
				old: "old",
				new: "",
			},
			expected: expected{
				isError: true,
			},
			mock: mock{
				status: 0,
				body:   "",
			},
		},
		{
			name: "empty old and new strings",
			fields: fields{
				Backlog: &backlog.Client{
					Writer:     io.Discard,
					BaseURL:    "https://example.com",
					APIKey:     "dummy",
					HTTPClient: &http.Client{},
				},
			},
			args: args{
				page: &Page{
					ID:   1,
					Name: "Old Name",
				},
				old: "",
				new: "",
			},
			expected: expected{
				isError: true,
			},
			mock: mock{
				status: 0,
				body:   "",
			},
		},
		{
			name: "api error",
			fields: fields{
				Backlog: &backlog.Client{
					Writer:     io.Discard,
					BaseURL:    "https://example.com",
					APIKey:     "dummy",
					HTTPClient: &http.Client{},
				},
			},
			args: args{
				page: &Page{
					ID:   1,
					Name: "Old Name",
				},
				old: "Old",
				new: "New",
			},
			expected: expected{
				isError: true,
			},
			mock: mock{
				status: 500,
				body:   `{"errors":[{"message":"Internal Server Error"}]}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Client{
				Client: tt.fields.Backlog,
			}
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			if tt.mock.status != 0 {
				httpmock.RegisterResponder(
					http.MethodPatch,
					fmt.Sprintf("%s/api/v2/wikis/%d?apiKey=%s", o.BaseURL, tt.args.page.ID, o.APIKey),
					httpmock.NewStringResponder(tt.mock.status, tt.mock.body),
				)
			}
			err := o.Rename(tt.args.page, tt.args.old, tt.args.new)
			if tt.expected.isError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestWiki_Replace(t *testing.T) {
	type fields struct {
		Backlog *backlog.Client
	}
	type args struct {
		page  *Page
		pairs []string
	}
	type expected struct {
		isError bool
	}
	type mock struct {
		status int
		body   string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		expected expected
		mock     mock
	}{
		{
			name: "basic",
			fields: fields{
				Backlog: &backlog.Client{
					Writer:     io.Discard,
					APIKey:     "dummy",
					HTTPClient: &http.Client{},
				},
			},
			args: args{
				page: &Page{
					ID:      1,
					Content: "Hello Old World",
				},
				pairs: []string{"Old", "New"},
			},
			expected: expected{
				isError: false,
			},
			mock: mock{
				status: 200,
				body:   "",
			},
		},
		{
			name: "empty page",
			fields: fields{
				Backlog: &backlog.Client{
					Writer:     io.Discard,
					APIKey:     "dummy",
					HTTPClient: &http.Client{},
				},
			},
			args: args{
				page:  nil,
				pairs: []string{"Old", "New"},
			},
			expected: expected{
				isError: true,
			},
			mock: mock{
				status: 0,
				body:   "",
			},
		},
		{
			name: "many pairs",
			fields: fields{
				Backlog: &backlog.Client{
					Writer:     io.Discard,
					APIKey:     "dummy",
					HTTPClient: &http.Client{},
				},
			},
			args: args{
				page: &Page{
					ID:      1,
					Content: "Hello Old World",
				},
				pairs: []string{"Old", "New", "World", "Earth"},
			},
			expected: expected{
				isError: true,
			},
			mock: mock{
				status: 0,
				body:   "",
			},
		},
		{
			name: "empty pairs",
			fields: fields{
				Backlog: &backlog.Client{
					Writer:     io.Discard,
					APIKey:     "dummy",
					HTTPClient: &http.Client{},
				},
			},
			args: args{
				page: &Page{
					ID:      1,
					Content: "Hello Old World",
				},
				pairs: nil,
			},
			expected: expected{
				isError: true,
			},
			mock: mock{
				status: 0,
				body:   "",
			},
		},
		{
			name: "invalid pairs",
			fields: fields{
				Backlog: &backlog.Client{
					Writer:     io.Discard,
					APIKey:     "dummy",
					HTTPClient: &http.Client{},
				},
			},
			args: args{
				page: &Page{
					ID:      1,
					Content: "Hello Old World",
				},
				pairs: []string{"Old"},
			},
			expected: expected{
				isError: true,
			},
			mock: mock{
				status: 0,
				body:   "",
			},
		},
		{
			name: "api error",
			fields: fields{
				Backlog: &backlog.Client{
					Writer:     io.Discard,
					APIKey:     "dummy",
					HTTPClient: &http.Client{},
				},
			},
			args: args{
				page: &Page{
					ID:      1,
					Content: "Hello Old World",
				},
				pairs: []string{"Old", "New"},
			},
			expected: expected{
				isError: true,
			},
			mock: mock{
				status: 500,
				body:   `{"errors":[{"message":"Internal Server Error"}]}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Client{
				Client: tt.fields.Backlog,
			}
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			if tt.mock.status != 0 {
				httpmock.RegisterResponder(
					http.MethodPatch,
					fmt.Sprintf("%s/api/v2/wikis/%d?apiKey=%s", o.BaseURL, tt.args.page.ID, o.APIKey),
					httpmock.NewStringResponder(tt.mock.status, tt.mock.body),
				)
			}
			err := o.Replace(tt.args.page, tt.args.pairs...)
			if tt.expected.isError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
