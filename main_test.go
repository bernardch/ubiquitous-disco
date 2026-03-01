package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// postCSV is our clean helper that mimics 'curl -F'.
// It returns the HTTP status code and the response body string.
func postCSV(t *testing.T, targetURL string, filePath string) (int, string) {
	file, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("Could not open %s: %v", filePath, err)
	}
	defer file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", filepath.Base(filePath))
	io.Copy(part, file)
	writer.Close()

	resp, err := http.Post(targetURL, writer.FormDataContentType(), body)
	if err != nil {
		t.Fatalf("Failed to make request to %s: %v", targetURL, err)
	}
	defer resp.Body.Close()

	responseBytes, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(responseBytes)
}

// setupTestServer spins up a real network server using your existing handlers.
// Note: This relies on your handlers being registered to http.DefaultServeMux,
// which happens automatically if you registered them with http.HandleFunc in main.go
// If they aren't registered during the test, we manually map them here.
func setupTestServer() *httptest.Server {
	mux := http.NewServeMux()

	// Map the endpoints to the handler logic (assuming you have your parseMatrix helper)
	// We wrap them exactly as they are in your main() function.
	mux.HandleFunc("/echo", echoHandler)
	mux.HandleFunc("/invert", invertHandler)
	mux.HandleFunc("/flatten", flattenHandler)
	mux.HandleFunc("/sum", sumHandler)
	mux.HandleFunc("/multiply", multiplyHandler)

	return httptest.NewServer(mux)
}

func TestAllMatrices(t *testing.T) {
	// 1. Start the test server
	server := setupTestServer()
	defer server.Close()

	// 2. Define all endpoints
	endpoints := []string{"/echo", "/invert", "/flatten", "/sum", "/multiply"}

	// 3. Define all your uploaded files and their expected behavior
	files := []struct {
		fileName      string
		isValid       bool   // true = expects 200 OK, false = expects 400 Bad Request
		expectedError string // substring to look for if invalid
	}{
		// Valid Files
		{"valid/valid-matrix.csv", true, ""},
		{"valid/valid2-matrix.csv", true, ""},
		{"valid/valid3-matrix.csv", true, ""},
		{"valid/matrix.csv", true, ""},
		{"valid/large-matrix.csv", true, ""},
		{"valid/whitespace-matrix.csv", true, ""}, // Should be valid because our parseMatrix trims spaces!

		// Invalid Files
		{"invalid/nonsquare-matrix.csv", false, "not square"},
		{"invalid/jagged-matrix.csv", false, "not square"},
		{"invalid/nonint-matrix.csv", false, "non-integer"},
	}

	// 4. Run the massive matrix of tests (Files * Endpoints)
	for _, f := range files {
		for _, endpoint := range endpoints {
			testName := f.fileName + endpoint

			t.Run(testName, func(t *testing.T) {
				// Send the request using our helper
				statusCode, response := postCSV(t, server.URL+endpoint, f.fileName)

				// Validate Status Codes
				if f.isValid {
					if statusCode != http.StatusOK {
						t.Errorf("[%s] Expected 200 OK, got %d. Response: %s", testName, statusCode, response)
					}
				} else {
					if statusCode != http.StatusBadRequest {
						t.Errorf("[%s] Expected 400 Bad Request, got %d. Response: %s", testName, statusCode, response)
					}
					// Validate the specific error message
					if !strings.Contains(strings.ToLower(response), f.expectedError) {
						t.Errorf("[%s] Expected error to contain '%s', got '%s'", testName, f.expectedError, response)
					}
				}
			})
		}
	}
}

// TestExactOutputs verifies that the core mathematical/formatting logic is perfectly accurate
// using the standard 3x3 matrix provided in the challenge requirements.
func TestExactOutputs(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	tests := []struct {
		endpoint string
		expected string
	}{
		{"/echo", "1,2,3\n4,5,6\n7,8,9\n"},
		{"/invert", "1,4,7\n2,5,8\n3,6,9\n"},
		{"/flatten", "1,2,3,4,5,6,7,8,9\n"},
		{"/sum", "45\n"},
		{"/multiply", "362880\n"},
	}

	for _, tt := range tests {
		t.Run("Standard Matrix "+tt.endpoint, func(t *testing.T) {
			statusCode, response := postCSV(t, server.URL+tt.endpoint, "matrix.csv")

			if statusCode != http.StatusOK {
				t.Fatalf("Expected 200 OK, got %d. Response: %s", statusCode, response)
			}

			// We use strings.TrimSpace on both to avoid failing over a trailing newline quirk
			if strings.TrimSpace(response) != strings.TrimSpace(tt.expected) {
				t.Errorf("Endpoint %s failed.\nExpected:\n%s\nGot:\n%s", tt.endpoint, tt.expected, response)
			}
		})
	}
}
