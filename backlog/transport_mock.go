package backlog

import (
	"bytes"
	"errors"
	"io"
	"net/http"
)

type mockRoundTripper struct {
	responses []*http.Response
	errors    []error
	n         int
}

func (o *mockRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	if o.n < len(o.errors) && o.errors[o.n] != nil {
		err := o.errors[o.n]
		o.n++
		return nil, err
	}
	if o.n < len(o.responses) {
		resp := o.responses[o.n]
		o.n++
		return resp, nil
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString("ok")),
		Header:     http.Header{},
	}, nil
}

type mockReadCloser struct{}

func (e *mockReadCloser) Read([]byte) (int, error) {
	return 0, errors.New("read error")
}

func (e *mockReadCloser) Close() error {
	return nil
}
