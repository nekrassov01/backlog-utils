package backlog

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// Backlog represents a Backlog client
type Backlog struct {
	Writer           io.Writer `json:"-"`
	BaseURL          string    `json:"baseUrl"`
	APIKey           string    `json:"-"`
	MaxRetryAttempts int       `json:"maxRetryAttempts"`
	MaxJitterMilli   int64     `json:"maxJitterMilli"`
}

// Do sends an HTTP request and retries if the response status code is 429 (Too Many Requests).
// It waits for the time specified in the X-RateLimit-Reset header before retrying.
func (o *Backlog) Do(req *http.Request) (*http.Response, error) {
	var n int
	var resp *http.Response
	var err error

	for {
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusTooManyRequests {
			return resp, nil
		}

		s := resp.Header.Get("X-RateLimit-Reset")
		wait := 1 * time.Second
		if s != "" {
			if ms, err := strconv.ParseInt(s, 10, 64); err == nil {
				now := time.Now().Unix()
				if ms > now {
					wait = time.Duration(ms-now) * time.Second
				}
			}
		}

		jitter := time.Duration(rand.Int63n(o.MaxJitterMilli)) * time.Millisecond // #nosec G404
		wait += jitter

		if _, err := io.Copy(io.Discard, resp.Body); err != nil {
			return nil, err
		}
		resp.Body.Close()

		time.Sleep(wait)
		n++

		if n > o.MaxRetryAttempts {
			return nil, fmt.Errorf("max retry attempts exceeded: %d", o.MaxRetryAttempts)
		}
	}
}
