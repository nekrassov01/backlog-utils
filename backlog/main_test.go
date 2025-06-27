package backlog

import (
	"io"
	"net/http"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	nowFunc = func() time.Time {
		return mustTime("2025-04-01T00:00:00Z")
	}
	defer func() {
		nowFunc = time.Now
	}()
	m.Run()
}

func mustTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

func newDefaultHTTPClient() *http.Client {
	return &http.Client{
		Transport: http.DefaultTransport,
	}
}

func newDefaultTransport() http.RoundTripper {
	return http.DefaultTransport
}

func mustNewRequest(method, url string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		panic(err)
	}
	return req
}
