package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	axios "github.com/MOHAMMADmiZAN/go-axios/axios"
)

// Post is a struct to parse JSON data for GET and POST requests
type Post struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func main() {
	// Initialize the Axios-like HTTP client with a global timeout of 15 seconds
	client := axios.NewClient(axios.Config{Timeout: 15}, nil)

	// Example of a GET request
	executeGetRequest(client)

	// Example of a POST request
	executePostRequest(client)

	// Example of a DELETE request
	executeDeleteRequest(client)

	// Example of a PUT request
	executePutRequest(client)
}

// GET Request Example
func executeGetRequest(client *axios.Client) {
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
		log.Printf("GET request failed: %v", err)
		return
	}

	// Parse the JSON response
	var post Post
	if err := json.Unmarshal(resp.Body, &post); err != nil {
		log.Printf("Failed to parse JSON response: %v", err)
		return
	}

	fmt.Printf("GET Response: %+v\n\n", post)
}

// POST Request Example
func executePostRequest(client *axios.Client) {
	// Create the JSON payload
	payload := Post{
		Title:  "foo",
		Body:   "bar",
		UserID: 1,
	}

	// Convert the payload to JSON
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal payload: %v", err)
		return
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
		log.Printf("POST request failed: %v", err)
		return
	}

	// Parse the JSON response
	var newPost Post
	if err := json.Unmarshal(resp.Body, &newPost); err != nil {
		log.Printf("Failed to parse JSON response: %v", err)
		return
	}

	fmt.Printf("POST Response: %+v\n\n", newPost)
}

// DELETE Request Example
func executeDeleteRequest(client *axios.Client) {
	// Configuration for the DELETE request
	reqConfig := axios.Config{
		Method: "DELETE",
		URL:    "https://jsonplaceholder.typicode.com/posts/1",
	}

	// Create a context with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Send the request
	resp, err := client.CancelableRequest(ctx, reqConfig)
	if err != nil {
		log.Printf("DELETE request failed: %v", err)
		return
	}

	fmt.Printf("DELETE Response Status: %s\n\n", resp.Status)
}

// PUT Request Example
func executePutRequest(client *axios.Client) {
	// Create the JSON payload for PUT
	payload := Post{
		ID:     1,
		Title:  "Updated Title",
		Body:   "Updated Body",
		UserID: 1,
	}

	// Convert the payload to JSON
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal payload: %v", err)
		return
	}

	// Configuration for the PUT request
	reqConfig := axios.Config{
		Method:  "PUT",
		URL:     "https://jsonplaceholder.typicode.com/posts/1",
		Body:    data,
		Headers: map[string][]string{"Content-Type": {"application/json"}},
	}

	// Create a context with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Send the request
	resp, err := client.CancelableRequest(ctx, reqConfig)
	if err != nil {
		log.Printf("PUT request failed: %v", err)
		return
	}

	// Parse the JSON response
	var updatedPost Post
	if err := json.Unmarshal(resp.Body, &updatedPost); err != nil {
		log.Printf("Failed to parse JSON response: %v", err)
		return
	}

	fmt.Printf("PUT Response: %+v\n\n", updatedPost)
}
