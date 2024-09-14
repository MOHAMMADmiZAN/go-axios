package axios

import (
	"fmt"
	"io"
	"net/http"
)

// RequestError represents an error that occurred during an HTTP request
type RequestError struct {
	StatusCode int
	Method     string
	URL        string
	Message    string
	Body       string // Optional: Store the response body for detailed error messages
}

// Error returns a detailed formatted error message
func (e *RequestError) Error() string {
	return fmt.Sprintf("request to %s %s failed with status code %d: %s\nResponse Body: %s",
		e.Method, e.URL, e.StatusCode, e.Message, e.Body)
}

// HandleResponseError creates a RequestError if the HTTP status code indicates an error
func HandleResponseError(resp *http.Response) error {
	if resp.StatusCode >= 400 {
		// Attempt to read the response body (optional for debugging)
		var responseBody string
		body, err := io.ReadAll(resp.Body)
		if err == nil && len(body) > 0 {
			responseBody = string(body)
		}

		// Return the error with status code and response details
		return &RequestError{
			StatusCode: resp.StatusCode,
			Method:     resp.Request.Method,
			URL:        resp.Request.URL.String(),
			Message:    http.StatusText(resp.StatusCode),
			Body:       responseBody,
		}
	}
	return nil
}
