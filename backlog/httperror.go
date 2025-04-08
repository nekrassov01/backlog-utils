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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	var errResp ErrorResponse
	if err := json.Unmarshal(body, &errResp); err != nil {
		return ""
	}

	messages := make([]string, 0, 2)
	for _, e := range errResp.Errors {
		messages = append(messages, e.Message)
	}

	if len(messages) == 0 {
		return ""
	}

	return strings.Join(messages, "; ")
}
