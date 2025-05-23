package backlog

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

// Error represents an error response from the Backlog API.
type Error struct {
	Message  string `json:"message"`
	Code     int    `json:"code"`
	MoreInfo string `json:"moreInfo"`
}

// ErrorResponse represents a response containing multiple errors from the Backlog API.
type ErrorResponse struct {
	Errors []Error `json:"errors"`
}

// GetErrorMessage reads the response body and returns the error message.
func GetErrorMessage(resp *http.Response) string {
	if resp == nil {
		return ""
	}

	s := ""
	if resp.StatusCode != 0 {
		s = http.StatusText(resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return s
	}

	var errResp ErrorResponse
	if err := json.Unmarshal(body, &errResp); err != nil {
		return s
	}

	messages := make([]string, 0, 2)
	for _, e := range errResp.Errors {
		messages = append(messages, e.Message)
	}

	if len(messages) == 0 {
		return s
	}

	return strings.Join(messages, "; ")
}
