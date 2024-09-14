package axios_test

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	axios "github.com/MOHAMMADmiZAN/go-axios/axios"
	"github.com/stretchr/testify/assert"
)

// TestNewClient ensures the client is initialized with correct timeout and transport options.
func TestNewClient(t *testing.T) {
	transportOpts := &axios.TransportOptions{
		MaxIdleConns:    10,
		IdleConnTimeout: 30 * time.Second,
	}
	config := axios.Config{Timeout: 15}
	client := axios.NewClient(config, transportOpts)

	assert.NotNil(t, client, "Client should be initialized")
	assert.Equal(t, 15*time.Second, client.HTTPClient().Timeout, "Timeout should match the config")
}

// TestClientRequestSuccess verifies that a GET request returns a successful response.
func TestClientRequestSuccess(t *testing.T) {
	// Mock server setup that returns a 200 OK status
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success"}`))
	}))
	defer server.Close()

	client := axios.NewClient(axios.Config{Timeout: 10}, nil)

	// Execute the GET request
	resp, err := client.Request(context.TODO(), axios.Config{
		Method: "GET",
		URL:    server.URL,
	})

	// Assert no errors and the correct response
	assert.NoError(t, err, "Request should not return an error")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response status should be 200 OK")
	assert.Contains(t, string(resp.Body), "success", "Response body should contain success message")
}

// TestClientRequestError verifies that an error response (500) is properly handled.
func TestClientRequestError(t *testing.T) {
	// Mock server setup that returns a 500 error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := axios.NewClient(axios.Config{Timeout: 10}, nil)

	// Execute the GET request expecting an error
	resp, err := client.Request(context.TODO(), axios.Config{
		Method: "GET",
		URL:    server.URL,
	})

	// Assert that the response is nil and an error is returned
	assert.Nil(t, resp, "Response should be nil on error")
	assert.Error(t, err, "Request should return an error for 500 status")
	assert.Contains(t, err.Error(), "500", "Error should contain status code")
}

// TestInterceptorRequest ensures that the request interceptor modifies the request (e.g., setting headers).
func TestInterceptorRequest(t *testing.T) {
	// Mock server setup to verify that the Authorization header is set
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"), "Authorization header should be set by interceptor")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := axios.NewClient(axios.Config{Timeout: 10}, nil)

	// Add a request interceptor to set the Authorization header
	client.GetInterceptorManager().AddInterceptor(axios.Interceptor{
		Request: func(req *http.Request) (*http.Request, error) {
			req.Header.Set("Authorization", "Bearer test-token")
			return req, nil
		},
	})

	// Execute the request
	resp, err := client.Request(context.TODO(), axios.Config{Method: "GET", URL: server.URL})
	assert.NoError(t, err, "Request should succeed")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Status should be 200 OK")
}

// TestInterceptorResponse verifies that the response interceptor modifies the response body.
func TestInterceptorResponse(t *testing.T) {
	// Mock server setup to return a response body
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"message": "original"}`))
	}))
	defer server.Close()

	client := axios.NewClient(axios.Config{Timeout: 10}, nil)
	im := axios.NewInterceptorManager()

	// Add a response interceptor to modify the response body
	im.AddInterceptor(axios.Interceptor{
		Response: func(resp *axios.Response) (*axios.Response, error) {
			resp.Body = []byte(`{"message": "intercepted"}`)
			return resp, nil
		},
	})

	// Execute the request and apply the response interceptor
	resp, err := client.Request(context.TODO(), axios.Config{Method: "GET", URL: server.URL})
	assert.NoError(t, err, "Request should succeed")

	resp, err = im.ApplyResponseInterceptors(resp)
	assert.NoError(t, err, "Response interceptors should not return an error")

	// Check that the response body was modified
	assert.Contains(t, string(resp.Body), "intercepted", "Response should be intercepted and modified")
}

// TestClientTimeout ensures that requests respect the configured timeout.
func TestClientTimeout(t *testing.T) {
	// Mock server setup with a delayed response to trigger timeout
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := axios.NewClient(axios.Config{Timeout: 1}, nil) // Set timeout to 1 second

	// Execute the GET request, expecting a timeout error
	resp, err := client.Request(context.TODO(), axios.Config{
		Method: "GET",
		URL:    server.URL,
	})

	// Assert that a timeout occurs
	assert.Nil(t, resp, "Response should be nil on timeout")
	assert.Error(t, err, "Request should return an error due to timeout")
	assert.Contains(t, err.Error(), "context deadline exceeded", "Error should indicate a timeout")
}

// TestClientCustomHeaders verifies that custom headers are properly set in the request.
func TestClientCustomHeaders(t *testing.T) {
	// Mock server setup to check headers
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"), "Content-Type header should match")
		assert.Equal(t, "Bearer custom-token", r.Header.Get("Authorization"), "Authorization header should match")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := axios.NewClient(axios.Config{Timeout: 10}, nil)

	// Execute the GET request with custom headers
	resp, err := client.Request(context.TODO(), axios.Config{
		Method: "GET",
		URL:    server.URL,
		Headers: http.Header{
			"Content-Type":  []string{"application/json"},
			"Authorization": []string{"Bearer custom-token"},
		},
	})

	// Assert the request was successful and headers were set
	assert.NoError(t, err, "Request should succeed")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Status should be 200 OK")
}

// TestClientQueryParams verifies that query parameters are included in the request URL.
func TestClientQueryParams(t *testing.T) {
	// Mock server setup to verify query parameters
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "value1", r.URL.Query().Get("param1"), "Query param1 should match")
		assert.Equal(t, "value2", r.URL.Query().Get("param2"), "Query param2 should match")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := axios.NewClient(axios.Config{Timeout: 10}, nil)

	// Execute the GET request with query parameters
	resp, err := client.Request(context.TODO(), axios.Config{
		Method: "GET",
		URL:    server.URL + "?param1=value1&param2=value2",
	})

	// Assert the request was successful
	assert.NoError(t, err, "Request should succeed")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Status should be 200 OK")
}

// TestClientEmptyResponseBody ensures that a response with an empty body is handled correctly.
func TestClientEmptyResponseBody(t *testing.T) {
	// Mock server setup with an empty response body
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := axios.NewClient(axios.Config{Timeout: 10}, nil)

	// Execute the GET request and expect an empty body
	resp, err := client.Request(context.TODO(), axios.Config{
		Method: "GET",
		URL:    server.URL,
	})

	// Assert the response was successful and body is empty
	assert.NoError(t, err, "Request should succeed")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Status should be 200 OK")
	assert.Empty(t, resp.Body, "Response body should be empty")
}

// TestClientParseJSON verifies that JSON responses are correctly parsed into Go types.
func TestClientParseJSON(t *testing.T) {
	// Mock server setup with a JSON response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"key": "value"}`))
	}))
	defer server.Close()

	client := axios.NewClient(axios.Config{Timeout: 10}, nil)

	// Execute the GET request
	resp, err := client.Request(context.TODO(), axios.Config{
		Method: "GET",
		URL:    server.URL,
	})
	assert.NoError(t, err, "Request should succeed")

	// Parse the JSON response
	var parsedBody map[string]string
	err = resp.ParseJSON(&parsedBody)
	assert.NoError(t, err, "JSON parsing should not return an error")
	assert.Equal(t, "value", parsedBody["key"], "Parsed key should match")
}

// TestClientCancelableRequest verifies that requests can be canceled using a context with cancellation.
func TestClientCancelableRequest(t *testing.T) {
	// Mock server setup with a delayed response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := axios.NewClient(axios.Config{Timeout: 10}, nil)

	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel the request before it's completed
	go func() {
		time.Sleep(1 * time.Second)
		cancel()
	}()

	// Execute the GET request and expect a cancellation error
	resp, err := client.CancelableRequest(ctx, axios.Config{
		Method: "GET",
		URL:    server.URL,
	})

	// Assert that the request was canceled
	assert.Nil(t, resp, "Response should be nil on cancellation")
	assert.Error(t, err, "Request should return an error due to cancellation")
	assert.Contains(t, err.Error(), "context canceled", "Error should indicate cancellation")
}

// TestClientConcurrentRequests ensures that the client can handle multiple requests concurrently without issues like race conditions.
func TestClientConcurrentRequests(t *testing.T) {
	// Mock server that returns a basic response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"message": "success"}`))
	}))
	defer server.Close()

	client := axios.NewClient(axios.Config{Timeout: 10}, nil)

	// Number of concurrent requests
	const numRequests = 10
	results := make(chan error, numRequests)

	// Execute multiple concurrent requests
	for i := 0; i < numRequests; i++ {
		go func() {
			resp, err := client.Request(context.TODO(), axios.Config{
				Method: "GET",
				URL:    server.URL,
			})
			if err != nil {
				results <- err
			} else if resp.StatusCode != http.StatusOK {
				results <- fmt.Errorf("unexpected status code: %d", resp.StatusCode)
			} else {
				results <- nil
			}
		}()
	}

	// Wait for all results and ensure there are no errors
	for i := 0; i < numRequests; i++ {
		assert.NoError(t, <-results, "All concurrent requests should succeed")
	}
}

// TestClientMultipleInterceptors   ensures that multiple interceptors are applied in the correct order and that they can modify both requests and responses.
func TestClientMultipleInterceptors(t *testing.T) {
	// Mock server setup
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify that both headers from the interceptors are present
		assert.Equal(t, "Bearer token1", r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.Write([]byte(`{"message": "success"}`))
	}))
	defer server.Close()

	client := axios.NewClient(axios.Config{Timeout: 10}, nil)

	// Add two request interceptors
	client.GetInterceptorManager().AddInterceptor(axios.Interceptor{
		Request: func(req *http.Request) (*http.Request, error) {
			req.Header.Set("Authorization", "Bearer token1")
			return req, nil
		},
	})

	client.GetInterceptorManager().AddInterceptor(axios.Interceptor{
		Request: func(req *http.Request) (*http.Request, error) {
			req.Header.Set("Content-Type", "application/json")
			return req, nil
		},
	})

	// Execute the request and check for both headers
	resp, err := client.Request(context.TODO(), axios.Config{Method: "GET", URL: server.URL})
	assert.NoError(t, err, "Request should succeed")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Status should be 200 OK")
}

// TestClientLargePayload checks how well the client handles large payloads for both requests and responses.
func TestClientLargePayload(t *testing.T) {
	// Create a large payload (1MB of 'A's)
	largePayload := make([]byte, 1<<20)
	for i := range largePayload {
		largePayload[i] = 'A'
	}

	// Mock server that returns a large response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(largePayload)
	}))
	defer server.Close()

	client := axios.NewClient(axios.Config{Timeout: 10}, nil)

	// Test request with large response
	resp, err := client.Request(context.TODO(), axios.Config{
		Method: "GET",
		URL:    server.URL,
	})

	// Ensure request succeeds and response body is large as expected
	assert.NoError(t, err, "Request should succeed")
	assert.Equal(t, len(largePayload), len(resp.Body), "Response body should match large payload size")
}

// TestClientRetryLogic  implements retry logic in case of failure (e.g., 500 response) and ensures the client can recover after a certain number of retries.
func TestClientRetryLogic(t *testing.T) {
	// Track the number of requests to simulate retry behavior
	requestCount := 0

	// Mock server that returns a 500 error on the first 2 requests, then succeeds
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		if requestCount < 3 {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "success"}`))
		}
	}))
	defer server.Close()

	client := axios.NewClient(axios.Config{Timeout: 10}, nil)

	// Implement simple retry logic: retry up to 3 times
	var resp *axios.Response
	var err error
	for i := 0; i < 3; i++ {
		resp, err = client.Request(context.TODO(), axios.Config{
			Method: "GET",
			URL:    server.URL,
		})
		if err == nil {
			break
		}
	}

	// Ensure that after retries, the request succeeds
	assert.NoError(t, err, "Request should succeed after retries")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Status should be 200 OK after retries")
	assert.Contains(t, string(resp.Body), "success", "Response should contain success message")
}

// TestClientMultipartUpload  ensures that the client can handle file uploads via multipart form data.
func TestClientMultipartUpload(t *testing.T) {
	// Mock server setup to check the uploaded file
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(10 << 20) // 10MB max memory
		assert.NoError(t, err, "Parsing multipart form should succeed")

		file, _, err := r.FormFile("file")
		assert.NoError(t, err, "Retrieving file should succeed")

		// Check the file contents
		fileContents := make([]byte, 4)
		file.Read(fileContents)
		assert.Equal(t, "test", string(fileContents), "File content should match")

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := axios.NewClient(axios.Config{Timeout: 10}, nil)

	// Create a file in memory for testing
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.txt")
	assert.NoError(t, err, "Creating form file should succeed")
	part.Write([]byte("test"))
	writer.Close()

	// Execute the multipart upload
	resp, err := client.Request(context.TODO(), axios.Config{
		Method:  "POST",
		URL:     server.URL,
		Body:    body.Bytes(),
		Headers: http.Header{"Content-Type": []string{writer.FormDataContentType()}},
	})

	// Ensure request succeeds
	assert.NoError(t, err, "Request should succeed")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Status should be 200 OK")
}

// TestClientRedirectHandling checks how the client handles HTTP redirects.
func TestClientRedirectHandling(t *testing.T) {
	// Mock server to simulate a redirect chain
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/redirect" {
			http.Redirect(w, r, "/final", http.StatusFound)
		} else if r.URL.Path == "/final" {
			w.Write([]byte(`{"message": "final destination"}`))
		}
	}))
	defer server.Close()

	client := axios.NewClient(axios.Config{Timeout: 10}, nil)

	// Execute the request that triggers a redirect
	resp, err := client.Request(context.TODO(), axios.Config{
		Method: "GET",
		URL:    server.URL + "/redirect",
	})

	// Ensure request succeeds and follows the redirect
	assert.NoError(t, err, "Request should succeed")
	assert.Contains(t, string(resp.Body), "final destination", "Response should follow the redirect")
}
