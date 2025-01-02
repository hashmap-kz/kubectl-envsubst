package cmd

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestReadRemote(t *testing.T) {

	t.Run("Successful Request", func(t *testing.T) {
		mockHTTPResponse := "Remote file content"
		http.DefaultClient = &http.Client{
			Transport: roundTripper(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(mockHTTPResponse)),
				}
			}),
		}
		result, err := readRemote("http://example.com/data")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if string(result) != mockHTTPResponse {
			t.Errorf("Expected: %s, Got: %s", mockHTTPResponse, string(result))
		}
	})

	t.Run("Failed Request", func(t *testing.T) {
		http.DefaultClient = &http.Client{
			Transport: roundTripper(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: 404,
					Body:       io.NopCloser(strings.NewReader("Not Found")),
				}
			}),
		}
		_, err := readRemote("http://example.com/not-found")
		if err == nil {
			t.Error("Expected error but got none")
		}
	})
}

// Helper for mocking http.Client
type roundTripper func(req *http.Request) *http.Response

func (f roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}
