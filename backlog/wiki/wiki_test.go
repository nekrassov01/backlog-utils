package wiki

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/nekrassov01/backlog-utils/backlog"
)

func TestNew(t *testing.T) {
	type args struct {
		url    string
		apiKey string
	}
	tests := []struct {
		name    string
		args    args
		want    *Wiki
		wantErr bool
	}{
		{
			name: "basic",
			args: args{
				url:    "https://example.com",
				apiKey: "dummy",
			},
			want: &Wiki{
				Backlog: &backlog.Backlog{
					Writer:           &bytes.Buffer{},
					BaseURL:          "https://example.com",
					APIKey:           "dummy",
					MaxRetryAttempts: 5,
					MaxJitterMilli:   1000,
				},
			},
			wantErr: false,
		},
		{
			name: "empty url",
			args: args{
				url:    "",
				apiKey: "dummy",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty api key",
			args: args{
				url:    "https://example.com",
				apiKey: "",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			got, err := New(w, tt.args.url, tt.args.apiKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWiki_List(t *testing.T) {
	type fields struct {
		Backlog *backlog.Backlog
	}
	type args struct {
		projectKey string
		pattern    string
	}
	type mock struct {
		status int
		body   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    mock
		want    []*Page
		wantErr bool
	}{
		{
			name: "basic",
			fields: fields{
				Backlog: &backlog.Backlog{
					Writer:           io.Discard,
					BaseURL:          "https://example.com",
					APIKey:           "dummy",
					MaxRetryAttempts: 2,
					MaxJitterMilli:   10,
				},
			},
			args: args{
				projectKey: "dummyProjectKey",
				pattern:    "",
			},
			mock: mock{
				status: 200,
				body:   `[{"id":1,"projectId":123,"name":"Test Page","content":""}]`,
			},
			want: []*Page{
				{
					ID:        1,
					ProjectID: 123,
					Name:      "Test Page",
					Content:   "",
				},
			},
			wantErr: false,
		},
		{
			name: "pattern",
			fields: fields{
				Backlog: &backlog.Backlog{
					Writer:           io.Discard,
					BaseURL:          "https://example.com",
					APIKey:           "dummy",
					MaxRetryAttempts: 2,
					MaxJitterMilli:   10,
				},
			},
			args: args{
				projectKey: "dummyProjectKey",
				pattern:    "Test",
			},
			mock: mock{
				status: 200,
				body:   `[{"id":1,"projectId":123,"name":"Test Page"},{"id":2,"projectId":123,"name":"Other Page"}]`,
			},
			want: []*Page{
				{
					ID:        1,
					ProjectID: 123,
					Name:      "Test Page",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid pattern",
			fields: fields{
				Backlog: &backlog.Backlog{
					Writer:           io.Discard,
					BaseURL:          "https://example.com",
					APIKey:           "dummy",
					MaxRetryAttempts: 2,
					MaxJitterMilli:   10,
				},
			},
			args: args{
				projectKey: "dummyProjectKey",
				pattern:    "[",
			},
			mock: mock{
				status: 0,
				body:   "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty url",
			fields: fields{
				Backlog: &backlog.Backlog{
					Writer:           io.Discard,
					BaseURL:          "",
					APIKey:           "dummy",
					MaxRetryAttempts: 2,
					MaxJitterMilli:   10,
				},
			},
			args: args{
				projectKey: "dummyProjectKey",
				pattern:    "",
			},
			mock: mock{
				status: 0,
				body:   "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty api key",
			fields: fields{
				Backlog: &backlog.Backlog{
					Writer:           io.Discard,
					BaseURL:          "https://example.com",
					APIKey:           "",
					MaxRetryAttempts: 2,
					MaxJitterMilli:   10,
				},
			},
			args: args{
				projectKey: "dummyProjectKey",
				pattern:    "",
			},
			mock: mock{
				status: 0,
				body:   "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty project key",
			fields: fields{
				Backlog: &backlog.Backlog{
					Writer:           io.Discard,
					BaseURL:          "https://example.com",
					APIKey:           "dummy",
					MaxRetryAttempts: 2,
					MaxJitterMilli:   10,
				},
			},
			args: args{
				projectKey: "",
				pattern:    "",
			},
			mock: mock{
				status: 0,
				body:   "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "api error",
			fields: fields{
				Backlog: &backlog.Backlog{
					Writer:           io.Discard,
					BaseURL:          "https://example.com",
					APIKey:           "dummy",
					MaxRetryAttempts: 2,
					MaxJitterMilli:   10,
				},
			},
			args: args{
				projectKey: "dummyProjectKey",
				pattern:    "",
			},
			mock: mock{
				status: 500,
				body:   `{"errors":[{"message":"Internal Server Error"}]}`,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid response",
			fields: fields{
				Backlog: &backlog.Backlog{
					Writer:           io.Discard,
					BaseURL:          "https://example.com",
					APIKey:           "dummy",
					MaxRetryAttempts: 2,
					MaxJitterMilli:   10,
				},
			},
			args: args{
				projectKey: "dummyProjectKey",
				pattern:    "",
			},
			mock: mock{
				status: 200,
				body:   `[{"id":1,"projectId":123,"name":}]`,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			if tt.mock.status != 0 {
				httpmock.RegisterResponder(
					http.MethodGet,
					fmt.Sprintf("%s/api/v2/wikis?projectIdOrKey=%s&apiKey=%s", tt.fields.Backlog.BaseURL, tt.args.projectKey, tt.fields.Backlog.APIKey),
					httpmock.NewStringResponder(tt.mock.status, tt.mock.body),
				)
			}
			o := &Wiki{
				Backlog: tt.fields.Backlog,
			}
			got, err := o.List(tt.args.projectKey, tt.args.pattern)
			if (err != nil) != tt.wantErr {
				t.Errorf("Wiki.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Wiki.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWiki_Get(t *testing.T) {
	type fields struct {
		Backlog *backlog.Backlog
	}
	type args struct {
		id int64
	}
	type mock struct {
		status int
		body   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    mock
		want    *Page
		wantErr bool
	}{
		{
			name: "basic",
			fields: fields{
				Backlog: &backlog.Backlog{
					Writer:  io.Discard,
					BaseURL: "https://example.com",
					APIKey:  "dummy",
				},
			},
			args: args{
				id: 1,
			},
			mock: mock{
				status: 200,
				body:   `{"id":1,"projectId":123,"name":"Test Page","content":"Sample Content"}`,
			},
			want: &Page{
				ID:        1,
				ProjectID: 123,
				Name:      "Test Page",
				Content:   "Sample Content",
			},
			wantErr: false,
		},
		{
			name: "empty url",
			fields: fields{
				Backlog: &backlog.Backlog{
					Writer:  io.Discard,
					BaseURL: "",
					APIKey:  "dummy",
				},
			},
			args: args{
				id: 1,
			},
			mock: mock{
				status: 0,
				body:   "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty api key",
			fields: fields{
				Backlog: &backlog.Backlog{
					Writer:  io.Discard,
					BaseURL: "https://example.com",
					APIKey:  "",
				},
			},
			args: args{
				id: 1,
			},
			mock: mock{
				status: 0,
				body:   "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid id",
			fields: fields{
				Backlog: &backlog.Backlog{
					Writer:  io.Discard,
					BaseURL: "https://example.com",
					APIKey:  "dummy",
				},
			},
			args: args{
				id: -1,
			},
			mock: mock{
				status: 0,
				body:   "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "api error",
			fields: fields{
				Backlog: &backlog.Backlog{
					Writer:  io.Discard,
					BaseURL: "https://example.com",
					APIKey:  "dummy",
				},
			},
			args: args{
				id: 1,
			},
			mock: mock{
				status: 404,
				body:   `{"errors":[{"message":"Not Found"}]}`,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid response",
			fields: fields{
				Backlog: &backlog.Backlog{
					Writer:  io.Discard,
					BaseURL: "https://example.com",
					APIKey:  "dummy",
				},
			},
			args: args{
				id: 1,
			},
			mock: mock{
				status: 200,
				body:   `{"id":1,"projectId":123,"name":}`,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			if tt.mock.status != 0 {
				httpmock.RegisterResponder(
					http.MethodGet,
					fmt.Sprintf("%s/api/v2/wikis/%d?apiKey=%s", tt.fields.Backlog.BaseURL, tt.args.id, tt.fields.Backlog.APIKey),
					httpmock.NewStringResponder(tt.mock.status, tt.mock.body),
				)
			}
			o := &Wiki{
				Backlog: tt.fields.Backlog,
			}
			got, err := o.Get(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Wiki.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Wiki.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWiki_Rename(t *testing.T) {
	type fields struct {
		Backlog *backlog.Backlog
	}
	type args struct {
		page *Page
		old  string
		new  string
	}
	type mock struct {
		status int
		body   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    mock
		wantErr bool
	}{
		{
			name: "basic",
			fields: fields{
				Backlog: &backlog.Backlog{
					Writer:  io.Discard,
					BaseURL: "https://example.com",
					APIKey:  "dummy",
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
			mock: mock{
				status: 200,
				body:   "",
			},
			wantErr: false,
		},
		{
			name: "empty page",
			fields: fields{
				Backlog: &backlog.Backlog{
					Writer:  io.Discard,
					BaseURL: "https://example.com",
					APIKey:  "dummy",
				},
			},
			args: args{
				page: nil,
				old:  "Old",
				new:  "New",
			},
			mock: mock{
				status: 0,
				body:   "",
			},
			wantErr: true,
		},
		{
			name: "empty old and new strings",
			fields: fields{
				Backlog: &backlog.Backlog{
					Writer:  io.Discard,
					BaseURL: "https://example.com",
					APIKey:  "dummy",
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
			mock: mock{
				status: 0,
				body:   "",
			},
			wantErr: true,
		},
		{
			name: "api error",
			fields: fields{
				Backlog: &backlog.Backlog{
					Writer:  io.Discard,
					BaseURL: "https://example.com",
					APIKey:  "dummy",
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
			mock: mock{
				status: 500,
				body:   `{"errors":[{"message":"Internal Server Error"}]}`,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			if tt.mock.status != 0 {
				httpmock.RegisterResponder(
					http.MethodPatch,
					fmt.Sprintf("%s/api/v2/wikis/%d?apiKey=%s", tt.fields.Backlog.BaseURL, tt.args.page.ID, tt.fields.Backlog.APIKey),
					httpmock.NewStringResponder(tt.mock.status, tt.mock.body),
				)
			}
			o := &Wiki{
				Backlog: tt.fields.Backlog,
			}
			err := o.Rename(tt.args.page, tt.args.old, tt.args.new)
			if (err != nil) != tt.wantErr {
				t.Errorf("Wiki.Rename() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWiki_Replace(t *testing.T) {
	type fields struct {
		Backlog *backlog.Backlog
	}
	type args struct {
		page  *Page
		pairs []string
	}
	type mock struct {
		status int
		body   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    mock
		wantErr bool
	}{
		{
			name: "basic",
			fields: fields{
				Backlog: &backlog.Backlog{
					Writer: io.Discard,
					APIKey: "dummy",
				},
			},
			args: args{
				page: &Page{
					ID:      1,
					Content: "Hello Old World",
				},
				pairs: []string{"Old", "New"},
			},
			mock: mock{
				status: 200,
				body:   "",
			},
			wantErr: false,
		},
		{
			name: "empty page",
			fields: fields{
				Backlog: &backlog.Backlog{
					Writer: io.Discard,
					APIKey: "dummy",
				},
			},
			args: args{
				page:  nil,
				pairs: []string{"Old", "New"},
			},
			mock: mock{
				status: 0,
				body:   "",
			},
			wantErr: true,
		},
		{
			name: "invalid pairs",
			fields: fields{
				Backlog: &backlog.Backlog{
					Writer: io.Discard,
					APIKey: "dummy",
				},
			},
			args: args{
				page: &Page{
					ID:      1,
					Content: "Hello Old World",
				},
				pairs: []string{"Old"},
			},
			mock: mock{
				status: 0,
				body:   "",
			},
			wantErr: true,
		},
		{
			name: "api error",
			fields: fields{
				Backlog: &backlog.Backlog{
					Writer: io.Discard,
					APIKey: "dummy",
				},
			},
			args: args{
				page: &Page{
					ID:      1,
					Content: "Hello Old World",
				},
				pairs: []string{"Old", "New"},
			},
			mock: mock{
				status: 500,
				body:   `{"errors":[{"message":"Internal Server Error"}]}`,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			if tt.mock.status != 0 {
				httpmock.RegisterResponder(
					http.MethodPatch,
					fmt.Sprintf("%s/api/v2/wikis/%d?apiKey=%s", tt.fields.Backlog.BaseURL, tt.args.page.ID, tt.fields.Backlog.APIKey),
					httpmock.NewStringResponder(tt.mock.status, tt.mock.body),
				)
			}
			o := &Wiki{
				Backlog: tt.fields.Backlog,
			}
			err := o.Replace(tt.args.page, tt.args.pairs...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Wiki.Replace() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
