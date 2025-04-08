package main

import (
	"context"
	"io"
	"testing"
)

func Test_cli(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "list empty url",
			args:    []string{name, "wiki", "list", "--base-url", "", "--api-key", "test", "--project-key", "test"},
			wantErr: true,
		},
		{
			name:    "list empty api key",
			args:    []string{name, "wiki", "list", "--base-url", "test", "--api-key", "", "--project-key", "test"},
			wantErr: true,
		},
		{
			name:    "list empty project key",
			args:    []string{name, "wiki", "list", "--base-url", "test", "--api-key", "test", "--project-key", ""},
			wantErr: true,
		},
		{
			name:    "list invalid pattern",
			args:    []string{name, "wiki", "list", "--base-url", "test", "--api-key", "test", "--project-key", "test", "--pattern", "["},
			wantErr: true,
		},
		{
			name:    "rename empty wiki id",
			args:    []string{name, "wiki", "rename", "--base-url", "test", "--api-key", "test", "--wiki-id", "", "--old", "old", "--new", "new"},
			wantErr: true,
		},
		{
			name:    "rename empty old string",
			args:    []string{name, "wiki", "rename", "--base-url", "test", "--api-key", "test", "--wiki-id", "test", "--old", "", "--new", "new"},
			wantErr: true,
		},
		{
			name:    "rename-all empty project key",
			args:    []string{name, "wiki", "rename-all", "--base-url", "test", "--api-key", "test", "--project-key", "", "--old", "old", "--new", "new"},
			wantErr: true,
		},
		{
			name:    "rename-all invalid pattern",
			args:    []string{name, "wiki", "rename-all", "--base-url", "test", "--api-key", "test", "--project-key", "test", "--pattern", "[", "--old", "old", "--new", "new"},
			wantErr: true,
		},
		{
			name:    "rename-all empty old string",
			args:    []string{name, "wiki", "rename-all", "--base-url", "test", "--api-key", "test", "--project-key", "test", "--pattern", "", "--old", "", "--new", "new"},
			wantErr: true,
		},
		{
			name:    "replace empty wiki id",
			args:    []string{name, "wiki", "replace", "--base-url", "test", "--api-key", "test", "--wiki-id", "", "--pairs", "key", "--pairs", "value"},
			wantErr: true,
		},
		{
			name:    "replace empty pairs",
			args:    []string{name, "wiki", "replace", "--base-url", "test", "--api-key", "test", "--wiki-id", "", "--pairs", "", "--pairs", ""},
			wantErr: true,
		},
		{
			name:    "replace empty invalid pairs",
			args:    []string{name, "wiki", "replace", "--base-url", "test", "--api-key", "test", "--wiki-id", "test", "--pairs", "key"},
			wantErr: true,
		},
		{
			name:    "replace-all empty project key",
			args:    []string{name, "wiki", "replace-all", "--base-url", "test", "--api-key", "test", "--project-key", "", "--pattern", "", "--pairs", "key", "--pairs", "value"},
			wantErr: true,
		},
		{
			name:    "replace-all invalid pattern",
			args:    []string{name, "wiki", "replace-all", "--base-url", "test", "--api-key", "test", "--project-key", "test", "--pattern", "[", "--pairs", "key", "--pairs", "value"},
			wantErr: true,
		},
		{
			name:    "replace-all empty pairs",
			args:    []string{name, "wiki", "replace-all", "--base-url", "test", "--api-key", "test", "--project-key", "test", "--pattern", "", "--pairs", "", "--pairs", ""},
			wantErr: true,
		},
		{
			name:    "replace-all invalid pairs",
			args:    []string{name, "wiki", "replace-all", "--base-url", "test", "--api-key", "test", "--project-key", "test", "--pattern", "", "--pairs", "key"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := newApp(io.Discard, io.Discard).Run(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
