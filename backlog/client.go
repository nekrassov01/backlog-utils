package backlog

import (
	"errors"
	"io"
	"net/http"
	"os"
	"time"
)

var nowFunc = time.Now

// Backlog represents a Backlog client.
type Client struct {
	BaseURL    string       `json:"baseUrl"`
	APIKey     string       `json:"-"`
	Writer     io.Writer    `json:"-"`
	HTTPClient *http.Client `json:"-"`
}

type ClientOption func(*Client)

func WithWriter(w io.Writer) ClientOption {
	if w == nil {
		w = os.Stdout
	}
	return func(o *Client) {
		o.Writer = w
	}
}

func WithTransport(transport http.RoundTripper) ClientOption {
	if transport == nil {
		transport = http.DefaultTransport
	}
	return func(o *Client) {
		o.HTTPClient.Transport = transport
	}
}

// NewClient creates a new Backlog client.
func NewClient(url, apiKey string, opts ...ClientOption) (*Client, error) {
	if url == "" {
		return nil, errors.New("empty URL")
	}
	if apiKey == "" {
		return nil, errors.New("empty API key")
	}
	o := &Client{
		BaseURL:    url,
		APIKey:     apiKey,
		Writer:     os.Stdout,
		HTTPClient: &http.Client{},
	}
	for _, opt := range opts {
		opt(o)
	}
	return o, nil
}
