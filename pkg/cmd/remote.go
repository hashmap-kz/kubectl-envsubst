package cmd

import (
	"fmt"
	"io"
	"net/http"
)

func readRemoteFileContent(url string) ([]byte, error) {

	// Make the HTTP GET request
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Check for HTTP errors
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cannot GET file content from: %s", url)
	}

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
