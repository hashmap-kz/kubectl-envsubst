package cmd

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func ReadRemoteFileContent(inputURL string) ([]byte, error) {
	// Parse and validate the URL
	parsedURL, err := url.Parse(inputURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return nil, fmt.Errorf("invalid URL: %s", inputURL)
	}

	// Make the HTTP GET request
	response, err := http.Get(parsedURL.String())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Check for HTTP errors
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cannot GET file content from: %s", inputURL)
	}

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
