package axios

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Response represents the parsed HTTP response
type Response struct {
	Status     string
	StatusCode int
	Body       []byte
	Headers    http.Header
}

// ParseResponse reads and parses the response body into a Response struct
func ParseResponse(resp *http.Response) (*Response, error) {
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	// Return the parsed response
	return &Response{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Body:       body,
		Headers:    resp.Header,
	}, nil
}

// ParseJSON parses the HTTP response body as JSON into the provided interface
func (r *Response) ParseJSON(v interface{}) error {
	if err := json.Unmarshal(r.Body, v); err != nil {
		return fmt.Errorf("error parsing JSON: %w", err)
	}
	return nil
}

// IsSuccess checks if the response has a 2xx status code
func (r *Response) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}
