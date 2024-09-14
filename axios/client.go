package axios

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// TransportOptions allows customization of http.Transport settings
type TransportOptions struct {
	MaxIdleConns        int
	IdleConnTimeout     time.Duration
	MaxIdleConnsPerHost int
	TLSHandshakeTimeout time.Duration
	ExpectContinue      time.Duration
}

// defaultTransport configures connection pooling and other transport settings
func defaultTransport(opts *TransportOptions) *http.Transport {
	if opts == nil {
		opts = &TransportOptions{
			MaxIdleConns:        100,
			IdleConnTimeout:     90 * time.Second,
			MaxIdleConnsPerHost: 100,
			TLSHandshakeTimeout: 10 * time.Second,
			ExpectContinue:      1 * time.Second,
		}
	}

	return &http.Transport{
		MaxIdleConns:          opts.MaxIdleConns,
		IdleConnTimeout:       opts.IdleConnTimeout,
		MaxIdleConnsPerHost:   opts.MaxIdleConnsPerHost,
		TLSHandshakeTimeout:   opts.TLSHandshakeTimeout,
		ExpectContinueTimeout: opts.ExpectContinue,
	}
}

// Client represents the HTTP client with custom configurations, transport, and interceptors
type Client struct {
	httpClient         *http.Client
	config             Config
	interceptorManager *InterceptorManager // Keep field unexported
}

// NewClient creates a new Client with a custom timeout and optional transport settings
func NewClient(config Config, transportOptions *TransportOptions) *Client {
	return &Client{
		httpClient: &http.Client{
			Transport: defaultTransport(transportOptions),
			Timeout:   time.Duration(config.Timeout) * time.Second,
		},
		config:             config,
		interceptorManager: NewInterceptorManager(),
	}
}

// GetInterceptorManager returns the interceptor manager for the client
func (c *Client) GetInterceptorManager() *InterceptorManager {
	return c.interceptorManager
}

// HTTPClient returns the internal http.Client (used for testing purposes)
func (c *Client) HTTPClient() *http.Client {
	return c.httpClient
}

// prepareRequestBody prepares the request body based on the config
func prepareRequestBody(config Config) (io.Reader, error) {
	if config.Body == nil {
		return nil, nil
	}
	return bytes.NewBuffer(config.Body), nil
}

// Request sends an HTTP request and returns the parsed response
func (c *Client) Request(ctx context.Context, config Config) (*Response, error) {
	finalConfig := mergeConfig(c.config, config)

	// Prepare the request body
	body, err := prepareRequestBody(finalConfig)
	if err != nil {
		return nil, fmt.Errorf("preparing request body: %w", err)
	}

	// Create a new request with context (supports timeout and cancellation)
	req, err := http.NewRequestWithContext(ctx, finalConfig.Method, finalConfig.URL, body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	// Apply request interceptors if any exist
	if c.interceptorManager != nil {
		req, err = c.interceptorManager.ApplyRequestInterceptors(req)
		if err != nil {
			return nil, fmt.Errorf("applying request interceptors: %w", err)
		}
	}

	// Set headers from config without overwriting existing ones
	for key, values := range finalConfig.Headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Execute the HTTP request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}

	// Check for HTTP errors (status code >= 400)
	if resp.StatusCode >= 400 {
		return nil, HandleResponseError(resp)
	}

	// Parse and return the response
	return ParseResponse(resp)
}

// CancelableRequest sends an HTTP request that supports cancellation via context
func (c *Client) CancelableRequest(ctx context.Context, config Config) (*Response, error) {
	return c.Request(ctx, config)
}
