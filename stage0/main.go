// Package main contains the first backend project in the learning roadmap.
package main

// Import the packages used by this HTTP learning server.
import (
	// Import context so shutdown can have a time limit.
	"context"
	// Import encoding/json for JSON request and response bodies.
	"encoding/json"
	// Import errors so handlers can compare known errors.
	"errors"
	// Import fmt for readable error messages.
	"fmt"
	// Import io for reading request bodies in the echo endpoint.
	"io"
	// Import log/slog for structured request logging.
	"log/slog"
	// Import net/http for the HTTP server, handlers, requests, and responses.
	"net/http"
	// Import os for environment variables and process output.
	"os"
	// Import strconv for converting URL path values from strings to integers.
	"strconv"
	// Import strings for trimming and validating text.
	"strings"
	// Import sync for protecting the in-memory quote slice.
	"sync"
	// Import time for timestamps and server timeouts.
	"time"
)

// Quote is the JSON representation of a quote stored by this service.
type Quote struct {
	// ID is the unique identifier returned to the client.
	ID int `json:"id"`
	// Text is the quote content supplied by the client.
	Text string `json:"text"`
	// Author is the person credited for the quote.
	Author string `json:"author"`
	// CreatedAt records when the quote was created.
	CreatedAt time.Time `json:"created_at"`
}

// QuoteStore is a small in-memory repository for learning HTTP before learning databases.
type QuoteStore struct {
	// mu prevents concurrent HTTP requests from changing the slice at the same time.
	mu sync.RWMutex
	// nextID is the identifier assigned to the next quote.
	nextID int
	// quotes contains the quotes currently held in memory.
	quotes []Quote
}

// newQuoteStore creates a store with a couple of records for GET requests.
func newQuoteStore() *QuoteStore {
	// Return a pointer so all handlers share the same store.
	return &QuoteStore{
		// Start IDs at one because zero is often used as an invalid ID.
		nextID: 1,
		// Seed data makes the GET endpoints useful immediately after startup.
		quotes: []Quote{
			// Add the first learning example.
			{ID: 1, Text: "The best way to learn is to build.", Author: "Unknown", CreatedAt: time.Now()},
			// Add the second learning example.
			{ID: 2, Text: "Clear code is better code.", Author: "Unknown", CreatedAt: time.Now()},
		},
	}
}

// add creates and stores a new quote.
func (s *QuoteStore) add(text string, author string) Quote {
	// Lock because this method changes shared state.
	s.mu.Lock()
	// Always unlock before returning, including future early returns.
	defer s.mu.Unlock()
	// Create a quote using the next available ID.
	quote := Quote{ID: s.nextID, Text: text, Author: author, CreatedAt: time.Now()}
	// Increase the ID so the next quote receives a different value.
	s.nextID++
	// Append the new quote to the in-memory collection.
	s.quotes = append(s.quotes, quote)
	// Return the created value to the handler.
	return quote
}

// list returns a copy so callers cannot accidentally modify the store without its lock.
func (s *QuoteStore) list() []Quote {
	// Acquire a read lock because this method only reads shared state.
	s.mu.RLock()
	// Always release the read lock before returning.
	defer s.mu.RUnlock()
	// Allocate a new slice with the same length as the stored collection.
	result := make([]Quote, len(s.quotes))
	// Copy the records into the new slice.
	copy(result, s.quotes)
	// Return the safe copy.
	return result
}

// find returns one quote or an error when the ID does not exist.
func (s *QuoteStore) find(id int) (Quote, error) {
	// Acquire a read lock because this method only reads shared state.
	s.mu.RLock()
	// Always release the read lock before returning.
	defer s.mu.RUnlock()
	// Search each stored quote.
	for _, quote := range s.quotes {
		// Return the matching quote immediately.
		if quote.ID == id {
			// Return the quote and no error.
			return quote, nil
		}
	}
	// Return a sentinel error so the handler can produce HTTP 404.
	return Quote{}, errors.New("quote not found")
}

// delete removes one quote and reports whether it was found.
func (s *QuoteStore) delete(id int) bool {
	// Lock because this method changes shared state.
	s.mu.Lock()
	// Always unlock before returning.
	defer s.mu.Unlock()
	// Search for the quote's position.
	for index, quote := range s.quotes {
		// Check whether this is the requested quote.
		if quote.ID == id {
			// Remove the matching item while preserving the other items.
			s.quotes = append(s.quotes[:index], s.quotes[index+1:]...)
			// Tell the caller that deletion succeeded.
			return true
		}
	}
	// Tell the caller that no quote had this ID.
	return false
}

// createQuoteRequest is the JSON body accepted by POST /quotes.
type createQuoteRequest struct {
	// Text is required and is decoded from the JSON "text" property.
	Text string `json:"text"`
	// Author is required and is decoded from the JSON "author" property.
	Author string `json:"author"`
}

// writeJSON sets JSON headers, writes a status code, and encodes a response body.
func writeJSON(w http.ResponseWriter, status int, value any) {
	// Tell the client that the response body contains UTF-8 JSON.
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// Write the HTTP status before writing the body.
	w.WriteHeader(status)
	// Encode the Go value as JSON and send it to the client.
	if err := json.NewEncoder(w).Encode(value); err != nil {
		// Log encoding failures because a response may already be partially written.
		slog.Error("encoding JSON response", "error", err)
	}
}

// writeError sends a consistent JSON error response.
func writeError(w http.ResponseWriter, status int, message string) {
	// Send an object instead of plain text so API clients receive one predictable shape.
	writeJSON(w, status, map[string]string{"error": message})
}

// methodNotAllowed tells a client which methods are accepted by a resource.
func methodNotAllowed(w http.ResponseWriter, allowed string) {
	// The Allow header documents the methods that can be used for this resource.
	w.Header().Set("Allow", allowed)
	// 405 means the method exists but is not supported for this resource.
	writeError(w, http.StatusMethodNotAllowed, "method not allowed")
}

// healthHandler demonstrates a simple successful response and a response header.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	// Add a custom header so clients can see how response headers work.
	w.Header().Set("X-Learning-Server", "stage0")
	// Return 200 because the server is available.
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// listQuotesHandler returns all quotes and reads an optional query parameter.
func listQuotesHandler(store *QuoteStore) http.HandlerFunc {
	// Return a function matching the http.HandlerFunc type.
	return func(w http.ResponseWriter, r *http.Request) {
		// Read the optional "author" value from /quotes?author=Unknown.
		authorFilter := strings.TrimSpace(r.URL.Query().Get("author"))
		// Read all current records from the store.
		quotes := store.list()
		// Create an empty result so JSON returns [] instead of null when there are no matches.
		filtered := make([]Quote, 0, len(quotes))
		// Inspect every quote returned by the repository.
		for _, quote := range quotes {
			// Skip quotes that do not match the optional author filter.
			if authorFilter != "" && !strings.EqualFold(quote.Author, authorFilter) {
				// Continue moves to the next quote.
				continue
			}
			// Add matching quotes to the response.
			filtered = append(filtered, quote)
		}
		// Return the collection with HTTP 200.
		writeJSON(w, http.StatusOK, filtered)
	}
}

// getQuoteHandler reads the {id} path parameter created by the ServeMux pattern.
func getQuoteHandler(store *QuoteStore) http.HandlerFunc {
	// Return a function matching the http.HandlerFunc type.
	return func(w http.ResponseWriter, r *http.Request) {
		// Read the path wildcard named "id" from /quotes/{id}.
		id, err := strconv.Atoi(r.PathValue("id"))
		// Reject values such as /quotes/abc with HTTP 400.
		if err != nil || id < 1 {
			// Explain that the path parameter must be a positive integer.
			writeError(w, http.StatusBadRequest, "id must be a positive integer")
			// Stop because the ID cannot be searched safely.
			return
		}
		// Search the store for the requested ID.
		quote, err := store.find(id)
		// Return 404 when the resource does not exist.
		if err != nil {
			// Tell the client that this specific resource was not found.
			writeError(w, http.StatusNotFound, "quote not found")
			// Stop after writing the error.
			return
		}
		// Return the matching quote with HTTP 200.
		writeJSON(w, http.StatusOK, quote)
	}
}

// createQuoteHandler decodes and validates a JSON request body.
func createQuoteHandler(store *QuoteStore) http.HandlerFunc {
	// Return a function matching the http.HandlerFunc type.
	return func(w http.ResponseWriter, r *http.Request) {
		// Limit the body so a learning server does not read an unbounded request.
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		// Create a value that will receive the JSON request.
		var input createQuoteRequest
		// Decode exactly one JSON object from the body.
		decoder := json.NewDecoder(r.Body)
		// Decode the first JSON value.
		if err := decoder.Decode(&input); err != nil {
			// Return 400 because malformed JSON is the client's input error.
			writeError(w, http.StatusBadRequest, "request body must be valid JSON")
			// Stop because there is no valid input to create.
			return
		}
		// Ensure a second JSON value was not appended after the first object.
		var extra any
		// Decode again and expect the end of the body.
		if err := decoder.Decode(&extra); err != io.EOF {
			// Reject multiple JSON values in one request body.
			writeError(w, http.StatusBadRequest, "request body must contain one JSON object")
			// Stop after reporting invalid input.
			return
		}
		// Trim whitespace before validating required text.
		input.Text = strings.TrimSpace(input.Text)
		// Trim whitespace before validating the author.
		input.Author = strings.TrimSpace(input.Author)
		// Reject missing required fields with HTTP 400.
		if input.Text == "" || input.Author == "" {
			// Explain which fields are required.
			writeError(w, http.StatusBadRequest, "text and author are required")
			// Stop because invalid data must not enter the store.
			return
		}
		// Save the validated input in memory.
		quote := store.add(input.Text, input.Author)
		// Location identifies the newly created resource.
		w.Header().Set("Location", fmt.Sprintf("/quotes/%d", quote.ID))
		// Return 201 because a new resource was created.
		writeJSON(w, http.StatusCreated, quote)
	}
}

// deleteQuoteHandler deletes a quote identified by its path parameter.
func deleteQuoteHandler(store *QuoteStore) http.HandlerFunc {
	// Return a function matching the http.HandlerFunc type.
	return func(w http.ResponseWriter, r *http.Request) {
		// Convert the string path value into an integer.
		id, err := strconv.Atoi(r.PathValue("id"))
		// Reject malformed or non-positive IDs.
		if err != nil || id < 1 {
			// Return HTTP 400 for invalid client input.
			writeError(w, http.StatusBadRequest, "id must be a positive integer")
			// Stop because the ID cannot be deleted.
			return
		}
		// Delete the requested quote.
		if !store.delete(id) {
			// Return HTTP 404 when there was no matching resource.
			writeError(w, http.StatusNotFound, "quote not found")
			// Stop after reporting the missing resource.
			return
		}
		// 204 means deletion succeeded and there is intentionally no response body.
		w.WriteHeader(http.StatusNoContent)
	}
}

// requestInfoHandler demonstrates request methods, headers, query values, cookies, and bodies.
func requestInfoHandler(w http.ResponseWriter, r *http.Request) {
	// Read the optional name query parameter.
	name := r.URL.Query().Get("name")
	// Read a client-provided header.
	clientHeader := r.Header.Get("X-Client-Name")
	// Read a cookie if one was sent.
	cookieValue := ""
	// Try to find the learning cookie.
	if cookie, err := r.Cookie("learning"); err == nil {
		// Copy the cookie value into the response object.
		cookieValue = cookie.Value
	}
	// Set a cookie so the browser can send it on a later request.
	http.SetCookie(w, &http.Cookie{Name: "learning", Value: "stage0", Path: "/"})
	// Return selected request information as JSON.
	writeJSON(w, http.StatusOK, map[string]any{
		// Include the HTTP method such as GET or POST.
		"method": r.Method,
		// Include the complete request URL.
		"url": r.URL.String(),
		// Include the query string value.
		"name_query": name,
		// Include the request header value.
		"client_header": clientHeader,
		// Include the received cookie value.
		"received_cookie": cookieValue,
		// Include the body size declared by the client when available.
		"content_length": r.ContentLength,
	})
}

// echoHandler reads a plain-text request body and sends it back.
func echoHandler(w http.ResponseWriter, r *http.Request) {
	// Limit this demonstration body to 64 KiB.
	r.Body = http.MaxBytesReader(w, r.Body, 64<<10)
	// Read all body bytes because this endpoint demonstrates body handling directly.
	body, err := io.ReadAll(r.Body)
	// Return 400 when the body cannot be read.
	if err != nil {
		// Explain that the request body was not readable.
		writeError(w, http.StatusBadRequest, "could not read request body")
		// Stop after reporting the invalid request.
		return
	}
	// Tell the client the response is plain text.
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	// Return 200 because the body was successfully echoed.
	w.WriteHeader(http.StatusOK)
	// Write the original bytes back to the client.
	_, _ = w.Write(body)
}

// formHandler parses application/x-www-form-urlencoded request bodies.
func formHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the form body and URL query values into r.Form.
	if err := r.ParseForm(); err != nil {
		// Return 400 when the form cannot be parsed.
		writeError(w, http.StatusBadRequest, "invalid form body")
		// Stop after reporting invalid input.
		return
	}
	// Return the parsed form values as JSON.
	writeJSON(w, http.StatusOK, map[string]any{
		// Form values can have multiple values for one key.
		"form": r.Form,
	})
}

// loggingMiddleware logs the method, path, status, and duration of every request.
func loggingMiddleware(next http.Handler) http.Handler {
	// Return a handler that wraps the next handler.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Record the start time before invoking the next handler.
		start := time.Now()
		// Call the wrapped handler so it can produce the response.
		next.ServeHTTP(w, r)
		// Log after the handler completes.
		slog.Info("http request", "method", r.Method, "path", r.URL.Path, "duration", time.Since(start))
	})
}

// buildHandler creates the complete route table for the service.
func buildHandler(store *QuoteStore) http.Handler {
	// Create the standard-library multiplexer.
	mux := http.NewServeMux()
	// Register an exact GET route for health checks.
	mux.HandleFunc("GET /health", healthHandler)
	// Register a GET collection route with an optional query parameter.
	mux.HandleFunc("GET /quotes", listQuotesHandler(store))
	// Register a GET route with a named path wildcard.
	mux.HandleFunc("GET /quotes/{id}", getQuoteHandler(store))
	// Register a POST route that accepts a JSON request body.
	mux.HandleFunc("POST /quotes", createQuoteHandler(store))
	// Register a DELETE route with a named path wildcard.
	mux.HandleFunc("DELETE /quotes/{id}", deleteQuoteHandler(store))
	// Register a GET route that demonstrates request metadata.
	mux.HandleFunc("GET /request-info", requestInfoHandler)
	// Register a POST route that echoes a plain-text request body.
	mux.HandleFunc("POST /echo", echoHandler)
	// Register a POST route that parses form-encoded request bodies.
	mux.HandleFunc("POST /form", formHandler)
	// Wrap routes with logging middleware before returning them.
	return loggingMiddleware(mux)
}

// main starts the server and waits for an interrupt before shutting down.
func main() {
	// Read the port from the environment, or use 8080 for local learning.
	port := os.Getenv("PORT")
	// Use a convenient default when PORT was not supplied.
	if port == "" {
		// Listen on port 8080.
		port = "8080"
	}
	// Create the complete handler tree.
	handler := buildHandler(newQuoteStore())
	// Configure the HTTP server and its basic timeouts.
	server := &http.Server{
		// Listen on all local interfaces at the selected port.
		Addr: ":" + port,
		// Attach the route handler.
		Handler: handler,
		// Limit time spent reading request headers.
		ReadHeaderTimeout: 5 * time.Second,
		// Limit time spent reading request bodies.
		ReadTimeout: 10 * time.Second,
		// Limit time spent writing responses.
		WriteTimeout: 10 * time.Second,
		// Close idle connections after this period.
		IdleTimeout: 60 * time.Second,
	}
	// Start the server in a goroutine so shutdown can be handled below.
	go func() {
		// Log the address before accepting requests.
		slog.Info("server listening", "address", server.Addr)
		// ListenAndServe returns http.ErrServerClosed during normal shutdown.
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			// Log unexpected startup or runtime errors.
			slog.Error("server stopped unexpectedly", "error", err)
			// Exit with a failure status for unexpected server errors.
			os.Exit(1)
		}
	}()
	// Wait for an operating-system signal such as Ctrl+C.
	waitForInterrupt()
	// Create a context that gives active requests five seconds to finish.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// Release the context timer after shutdown completes.
	defer cancel()
	// Ask the server to stop accepting new requests and finish active work.
	if err := server.Shutdown(ctx); err != nil {
		// Log shutdown failures explicitly.
		slog.Error("server shutdown failed", "error", err)
		// Exit with a failure status when graceful shutdown did not finish.
		os.Exit(1)
	}
	// Tell the learner that shutdown completed.
	slog.Info("server stopped")
}

// waitForInterrupt blocks until the user presses Ctrl+C on supported systems.
func waitForInterrupt() {
	// Create a channel that receives operating-system signals.
	signals := make(chan os.Signal, 1)
	// Register interrupts such as Ctrl+C with the channel.
	signalNotify(signals)
	// Wait until a signal arrives.
	<-signals
}
