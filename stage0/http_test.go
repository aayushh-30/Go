// Package main tests the HTTP behavior exposed by the stage0 server.
package main

// Import the packages required to create requests and inspect responses.
import (
	// Import bytes to provide a request body from a byte slice.
	"bytes"
	// Import encoding/json to decode JSON responses.
	"encoding/json"
	// Import net/http for HTTP constants and client behavior.
	"net/http"
	// Import net/http/httptest for an in-memory HTTP server.
	"net/http/httptest"
	// Import strconv to convert the created ID into a URL path.
	"strconv"
	// Import testing for Go's test runner.
	"testing"
)

// TestCreateAndReadQuote verifies a JSON POST followed by a JSON GET.
func TestCreateAndReadQuote(t *testing.T) {
	// Build the handler without opening a network port.
	handler := buildHandler(newQuoteStore())
	// Create a JSON request body for the POST endpoint.
	body := bytes.NewBufferString(`{"text":"Practice makes progress.","author":"Learner"}`)
	// Create a POST request with the JSON body.
	request := httptest.NewRequest(http.MethodPost, "/quotes", body)
	// Tell the server how to interpret the request body.
	request.Header.Set("Content-Type", "application/json")
	// Record the response produced by the handler.
	response := httptest.NewRecorder()
	// Execute the request entirely in memory.
	handler.ServeHTTP(response, request)
	// Check that creating a resource returned 201 Created.
	if response.Code != http.StatusCreated {
		// Stop the test with the unexpected status code.
		t.Fatalf("expected status %d, got %d", http.StatusCreated, response.Code)
	}
	// Decode the response JSON into a Quote value.
	var created Quote
	// Fail if the response is not valid JSON.
	if err := json.NewDecoder(response.Body).Decode(&created); err != nil {
		// Stop the test with the decoding error.
		t.Fatalf("decode created quote: %v", err)
	}
	// Check that the server returned the submitted text.
	if created.Text != "Practice makes progress." {
		// Report the unexpected response value.
		t.Fatalf("expected created text, got %q", created.Text)
	}
	// Build a GET request using the returned resource ID.
	getRequest := httptest.NewRequest(http.MethodGet, "/quotes/"+itoa(created.ID), nil)
	// Record the GET response.
	getResponse := httptest.NewRecorder()
	// Execute the GET request.
	handler.ServeHTTP(getResponse, getRequest)
	// Check that reading the created resource returned 200 OK.
	if getResponse.Code != http.StatusOK {
		// Stop the test with the unexpected status code.
		t.Fatalf("expected status %d, got %d", http.StatusOK, getResponse.Code)
	}
}

// TestInvalidJSONReturnsBadRequest verifies validation of malformed input.
func TestInvalidJSONReturnsBadRequest(t *testing.T) {
	// Build the handler with a fresh store.
	handler := buildHandler(newQuoteStore())
	// Create a request containing malformed JSON.
	request := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewBufferString(`{"text":`))
	// Mark the request body as JSON.
	request.Header.Set("Content-Type", "application/json")
	// Record the response.
	response := httptest.NewRecorder()
	// Execute the invalid request.
	handler.ServeHTTP(response, request)
	// Confirm that invalid client input receives 400 Bad Request.
	if response.Code != http.StatusBadRequest {
		// Report the unexpected status code.
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.Code)
	}
}

// TestMethodNotAllowed verifies that an unsupported method receives 405.
func TestMethodNotAllowed(t *testing.T) {
	// Build the handler with a fresh store.
	handler := buildHandler(newQuoteStore())
	// Create a PATCH request even though /health only supports GET.
	request := httptest.NewRequest(http.MethodPatch, "/health", nil)
	// Record the response.
	response := httptest.NewRecorder()
	// Execute the unsupported request.
	handler.ServeHTTP(response, request)
	// Confirm that the standard mux rejects the method.
	if response.Code != http.StatusMethodNotAllowed {
		// Report the unexpected status code.
		t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, response.Code)
	}
}

// itoa avoids adding another package to this small test explanation.
func itoa(value int) string {
	// Convert the integer ID to its decimal string representation.
	return strconv.Itoa(value)
}
