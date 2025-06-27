package backlog

import (
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"
)

const resetHeaderKey = "X-Ratelimit-Reset"

var _ http.RoundTripper = (*RetryableTransport)(nil)

var retryableStatus = map[int]struct{}{
	http.StatusTooManyRequests:     {}, // 429
	http.StatusInternalServerError: {}, // 500
	http.StatusBadGateway:          {}, // 502
	http.StatusServiceUnavailable:  {}, // 503
	http.StatusGatewayTimeout:      {}, // 504
}

// RetryableTransport is a custom HTTP transport that retries requests on certain status codes.
type RetryableTransport struct {
	Transport        http.RoundTripper `json:"-"`
	InitialInterval  time.Duration     `json:"initialInterval"`
	MaxInterval      time.Duration     `json:"maxInterval"`
	MaxRetryAttempts int               `json:"maxRetryAttempts"`
	MaxJitterMilli   int               `json:"maxJitterMilli"`
}

// NewRetryableTransport creates a new Transport with the specified parameters.
func NewRetryableTransport(initialInterval, maxInterval time.Duration, maxRetryAttempts, maxJitterMilli int) *RetryableTransport {
	if initialInterval < 0 {
		initialInterval = 1 * time.Second
	}
	if maxInterval < 0 {
		maxInterval = 30 * time.Second
	}
	if maxRetryAttempts < 0 {
		maxRetryAttempts = 5
	}
	if maxJitterMilli < 0 {
		maxJitterMilli = 3000
	}
	t := &RetryableTransport{
		Transport:        http.DefaultTransport,
		InitialInterval:  initialInterval,
		MaxInterval:      maxInterval,
		MaxRetryAttempts: maxRetryAttempts,
		MaxJitterMilli:   maxJitterMilli,
	}
	return t
}

// RoundTrip sends an HTTP request and retries if the response status code is 429 (Too Many Requests).
// It waits for the time specified in the X-RateLimit-Reset header before retrying.
// The retry attempts are limited by MaxRetryAttempts and a random jitter is added to the wait time.
// The jitter is a random duration between 0 and MaxJitterMilli milliseconds.
func (o *RetryableTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Body != nil && req.GetBody == nil {
		return nil, errors.New("request body is not rewindable")
	}

	status := 0

	for i := range o.MaxRetryAttempts {
		var err error
		if req.Body != nil && req.GetBody != nil && i > 0 {
			req.Body, err = req.GetBody()
			if err != nil {
				return nil, fmt.Errorf("failed to rewind request body: %w", err)
			}
		}
		resp, err := o.Transport.RoundTrip(req)
		if err != nil {
			return nil, fmt.Errorf("failed to request: %w", err)
		}
		status = resp.StatusCode

		if _, ok := retryableStatus[status]; !ok {
			return resp, nil
		}

		if _, err := io.Copy(io.Discard, resp.Body); err != nil {
			err2 := fmt.Errorf("failed to read response body: %w", err)
			if err3 := resp.Body.Close(); err3 != nil {
				err4 := errors.Join(err2, fmt.Errorf("failed to close response body: %w", err3))
				return nil, err4
			}
			return nil, err2
		}
		if err := resp.Body.Close(); err != nil {
			err2 := fmt.Errorf("failed to close response body: %w", err)
			return nil, err2
		}

		interval := o.InitialInterval
		if resp.StatusCode == http.StatusTooManyRequests {
			s := resp.Header.Get(resetHeaderKey)
			if s != "" {
				if seconds, err := strconv.ParseInt(s, 10, 64); err == nil {
					now := time.Now().Unix()
					if seconds > now {
						interval = time.Duration(seconds-now) * time.Second
					}
				}
			}
		} else {
			interval = min(o.InitialInterval<<i, o.MaxInterval)
		}

		jitter := time.Duration(rand.N(o.MaxJitterMilli)) * time.Millisecond // #nosec G404
		interval += jitter
		timer := time.NewTimer(interval)

		select {
		case <-timer.C:
		case <-req.Context().Done():
			timer.Stop()
			return nil, req.Context().Err()
		}

		timer.Stop()
	}

	return nil, errors.New("max retry attempts exceeded")
}
