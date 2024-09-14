
# Go-Axios

Go-Axios is a Golang HTTP client inspired by Axios JS. It's not an official package but aims to bring a similar approach to handling HTTP requests in Go. The package provides support for features like request/response interceptors, customizable timeouts, JSON handling, and request cancellation using contexts, offering a simple and efficient way to manage HTTP requests.

## Features
- **Simple HTTP Requests**: Supports `GET`, `POST`, `PUT`, `DELETE`, `PATCH`, and other HTTP methods.
- **Customizable Timeouts**: Allows configuring global and per-request timeouts.
- **Request & Response Interceptors**: Middleware to modify requests and responses.
- **JSON Handling**: Easy-to-use JSON serialization and deserialization.
- **Context-Aware Requests**: Supports request cancellation and deadlines via `context.Context`.
- **Custom Transport Options**: Supports configuration of connection pooling, TLS settings, and more.
- **Error Handling**: Detailed error messages with request method, URL, status codes, and response body.

---

## Installation

To install the `go-axios` package, run:

```bash
go get github.com/MOHAMMADmiZAN/go-axios/axios@v1.1.0
```

---

## Usage

### Example 1: Basic GET Request

This example shows how to perform a basic GET request using `go-axios`:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    axios "github.com/MOHAMMADmiZAN/go-axios/axios"
)

func main() {
    // Initialize the Axios-like HTTP client with a global timeout of 15 seconds
    client := axios.NewClient(axios.Config{Timeout: 15}, nil)

    // Configuration for the GET request
    reqConfig := axios.Config{
        Method: "GET",
        URL:    "https://jsonplaceholder.typicode.com/posts/1",
    }

    // Create a context with a timeout of 5 seconds
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Send the request
    resp, err := client.CancelableRequest(ctx, reqConfig)
    if err != nil {
        log.Fatalf("Request failed: %v", err)
    }

    fmt.Printf("Response: %s\n", string(resp.Body))
}
```

### Example 2: POST Request with JSON Payload

This example shows how to send a `POST` request with a JSON payload:

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "time"

    axios "github.com/MOHAMMADmiZAN/go-axios/axios"
)

// Struct to represent the data sent in the POST request
type Post struct {
    UserID int    `json:"userId"`
    Title  string `json:"title"`
    Body   string `json:"body"`
}

func main() {
    // Initialize the Axios-like HTTP client
    client := axios.NewClient(axios.Config{Timeout: 15}, nil)

    // Create a JSON payload
    postData := Post{
        UserID: 1,
        Title:  "foo",
        Body:   "bar",
    }

    // Convert the payload to JSON
    data, err := json.Marshal(postData)
    if err != nil {
        log.Fatalf("Error marshaling JSON: %v", err)
    }

    // Configuration for the POST request
    reqConfig := axios.Config{
        Method:  "POST",
        URL:     "https://jsonplaceholder.typicode.com/posts",
        Body:    data,
        Headers: map[string][]string{"Content-Type": {"application/json"}},
    }

    // Create a context with a timeout of 5 seconds
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Send the request
    resp, err := client.CancelableRequest(ctx, reqConfig)
    if err != nil {
        log.Fatalf("Request failed: %v", err)
    }

    fmt.Printf("Response: %s\n", string(resp.Body))
}
```

---

## Features in Detail

### 1. **Custom Timeouts**
   - You can specify a global timeout for the client, and also override it per request using the `context.WithTimeout` function.

### 2. **Request & Response Interceptors**
   - Interceptors allow you to modify requests and responses, which is useful for adding authentication tokens, logging, or modifying data before/after the request is sent.

   ```go
   client := axios.NewClient(axios.Config{Timeout: 15}, nil)

   // Add a request interceptor
   client.AddInterceptor(axios.Interceptor{
       Request: func(req *http.Request) (*http.Request, error) {
           // Example: Adding Authorization header to every request
           req.Header.Set("Authorization", "Bearer YOUR_TOKEN_HERE")
           return req, nil
       },
   })

   // Add a response interceptor
   client.AddInterceptor(axios.Interceptor{
       Response: func(resp *axios.Response) (*axios.Response, error) {
           // Example: Logging the response status
           fmt.Printf("Response Status: %s\n", resp.Status)
           return resp, nil
       },
   })
   ```

### 3. **Custom Transport Options**
   - Customize the underlying transport options like connection pooling, TLS settings, and more:

   ```go
   transportOptions := &axios.TransportOptions{
       MaxIdleConns:        100,
       IdleConnTimeout:     90 * time.Second,
       MaxIdleConnsPerHost: 10,
       TLSHandshakeTimeout: 10 * time.Second,
   }

   client := axios.NewClient(axios.Config{Timeout: 15}, transportOptions)
   ```

### 4. **Error Handling**
   - The `go-axios` package provides detailed error messages that include the request method, URL, status code, and an optional response body for easier debugging.

   ```go
   resp, err := client.CancelableRequest(ctx, reqConfig)
   if err != nil {
       if reqErr, ok := err.(*axios.RequestError); ok {
           log.Printf("Request failed: %s %s with status code %d\nResponse: %s\n",
               reqErr.Method, reqErr.URL, reqErr.StatusCode, reqErr.Body)
       }
   }
   ```

### 5. **Context-Aware Requests**
   - You can cancel requests or set deadlines using the standard `context.Context` in Go:

   ```go
   ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
   defer cancel()

   resp, err := client.CancelableRequest(ctx, reqConfig)
   if err != nil {
       log.Fatalf("Request failed: %v", err)
   }
   ```

---

## Configuration

The `Config` struct defines the client configuration:

```go
type Config struct {
    Method  string
    URL     string
    Headers http.Header
    Params  map[string]string
    Body    []byte
    Timeout int
}
```

- `Method`: HTTP method (`GET`, `POST`, `PUT`, etc.).
- `URL`: The endpoint URL.
- `Headers`: Optional HTTP headers.
- `Params`: Optional query parameters.
- `Body`: Optional body data (for `POST`, `PUT`, etc.).
- `Timeout`: Request timeout in seconds (overridden by context timeouts).

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request. For major changes, please open an issue first to discuss what you would like to change.

---

## Conclusion

`go-axios` makes working with HTTP requests in Golang simple and efficient, while still offering the flexibility to handle complex scenarios such as custom timeouts, interceptors, and detailed error reporting.

For more advanced usage and customization, feel free to check out the examples in this README or explore the source code in the repository.

---

### To Summarize:
- **Install**: `go get github.com/MOHAMMADmiZAN/go-axios/axios`
- **Simple HTTP requests**: Easily perform GET, POST, PUT, DELETE, etc.
- **Customizable**: Support for timeouts, interceptors, and custom transport.
- **Error Handling**: Detailed error reporting to assist debugging.

We hope this library helps you manage HTTP requests with ease in your Go projects!
