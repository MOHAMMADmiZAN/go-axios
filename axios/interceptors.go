package axios

import (
	"fmt"
	"net/http"
)

// Interceptor defines functions for request and response interception
type Interceptor struct {
	Request  func(*http.Request) (*http.Request, error)
	Response func(*Response) (*Response, error)
}

// InterceptorManager manages the addition and execution of interceptors
type InterceptorManager struct {
	interceptors []Interceptor
}

// NewInterceptorManager initializes a new InterceptorManager
func NewInterceptorManager() *InterceptorManager {
	return &InterceptorManager{
		interceptors: nil, // Lazy initialization, avoids allocating an empty slice unnecessarily
	}
}

// AddInterceptor registers an interceptor, lazily initializing the slice if needed
func (im *InterceptorManager) AddInterceptor(i Interceptor) {
	if im.interceptors == nil {
		im.interceptors = []Interceptor{}
	}
	im.interceptors = append(im.interceptors, i)
}

// ApplyRequestInterceptors applies all request interceptors in sequence, stopping if any returns an error
func (im *InterceptorManager) ApplyRequestInterceptors(req *http.Request) (*http.Request, error) {
	var err error
	for idx, interceptor := range im.interceptors {
		req, err = interceptor.Request(req)
		if err != nil {
			return nil, fmt.Errorf("request interceptor %d failed: %w", idx, err)
		}
	}
	return req, nil
}

// ApplyResponseInterceptors applies all response interceptors in sequence, stopping if any returns an error
func (im *InterceptorManager) ApplyResponseInterceptors(resp *Response) (*Response, error) {
	var err error
	for idx, interceptor := range im.interceptors {
		resp, err = interceptor.Response(resp)
		if err != nil {
			return nil, fmt.Errorf("response interceptor %d failed: %w", idx, err)
		}
	}
	return resp, nil
}
