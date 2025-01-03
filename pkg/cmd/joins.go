package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

func JoinFiles(files []string, hasStdin bool) ([]byte, error) {
	buf := bytes.Buffer{}

	totalFiles := len(files)
	if hasStdin {
		totalFiles += 1
	}
	needSeparator := totalFiles > 1
	const separator = "\n---\n"

	// process STDIN
	if hasStdin {
		stdin, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, err
		}
		buf.Write(stdin)
		if needSeparator {
			buf.WriteString(separator)
		}
	}

	// process files
	for _, f := range files {
		// passed as URL
		if IsURL(f) {
			data, err := readRemote(f)
			if err != nil {
				return nil, err
			}
			buf.Write(data)
			if needSeparator {
				buf.WriteString(separator)
			}
			continue
		}

		// plain file
		file, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}
		buf.Write(file)
		if needSeparator {
			buf.WriteString(separator)
		}
	}

	return buf.Bytes(), nil
}

func readRemote(url string) ([]byte, error) {

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
