package axios

import "net/http"

// Config stores the HTTP request configuration options
type Config struct {
	Method  string
	URL     string
	Headers http.Header
	Params  map[string]string
	Body    []byte
	Timeout int
}

// mergeConfig merges default and user-defined configurations
func mergeConfig(defaultConfig, userConfig Config) Config {
	finalConfig := defaultConfig

	// Merge HTTP method
	if userConfig.Method != "" {
		finalConfig.Method = userConfig.Method
	}

	// Merge URL
	if userConfig.URL != "" {
		finalConfig.URL = userConfig.URL
	}

	// Merge Headers
	finalConfig.Headers = mergeHeaders(defaultConfig.Headers, userConfig.Headers)

	// Merge Query Params
	finalConfig.Params = mergeParams(defaultConfig.Params, userConfig.Params)

	// Merge Body
	if userConfig.Body != nil {
		finalConfig.Body = userConfig.Body
	}

	// Merge Timeout
	if userConfig.Timeout != 0 {
		finalConfig.Timeout = userConfig.Timeout
	}

	return finalConfig
}

// mergeHeaders merges two HTTP header sets with user-defined headers overriding the defaults
func mergeHeaders(defaultHeaders, userHeaders http.Header) http.Header {
	if defaultHeaders == nil {
		defaultHeaders = http.Header{}
	}

	for key, values := range userHeaders {
		for _, value := range values {
			defaultHeaders.Set(key, value) // Overwrites existing headers
		}
	}

	return defaultHeaders
}

// mergeParams merges query parameters, prioritizing user-defined ones
func mergeParams(defaultParams, userParams map[string]string) map[string]string {
	if defaultParams == nil {
		defaultParams = make(map[string]string)
	}

	for key, value := range userParams {
		defaultParams[key] = value // Overwrites existing parameters
	}

	return defaultParams
}
