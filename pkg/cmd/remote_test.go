package cmd

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReadRemoteFileContent(t *testing.T) {
	tests := []struct {
		name         string
		inputURL     string
		mockResponse string
		statusCode   int
		wantError    bool
	}{
		{
			name:         "Valid URL and HTTP response",
			inputURL:     "http://example.com/test",
			mockResponse: "file content",
			statusCode:   http.StatusOK,
			wantError:    false,
		},
		{
			name:      "Invalid URL - malformed",
			inputURL:  ":invalid-url",
			wantError: true,
		},
		{
			name:      "Invalid URL - missing host",
			inputURL:  "http:///test",
			wantError: true,
		},
		{
			name:       "Non-200 HTTP status code",
			inputURL:   "http://example.com/notfound",
			statusCode: http.StatusNotFound,
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock HTTP server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.statusCode)
				fmt.Fprint(w, tt.mockResponse)
			}))
			defer server.Close()

			// Replace the input URL with the mock server's URL if statusCode is set
			inputURL := tt.inputURL
			if tt.statusCode != 0 {
				inputURL = server.URL
			}

			// Call the function under test
			got, err := ReadRemoteFileContent(inputURL)

			// Check for errors
			if tt.wantError {
				if err == nil {
					t.Errorf("Expected an error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error but got one: %v", err)
				}
				// Verify the response content
				if string(got) != tt.mockResponse {
					t.Errorf("Expected response %q, got %q", tt.mockResponse, got)
				}
			}
		})
	}
}
