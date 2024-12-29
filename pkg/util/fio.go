package util

import (
	"fmt"
	"io"
	"os"
)

func ReadFromStdin() ([]byte, error) {

	// Check if stdin is piped or redirected
	info, err := os.Stdin.Stat()
	if err != nil {
		return nil, fmt.Errorf("error checking stdin: %w", err)
	}

	// If stdin is a terminal (no pipe/file), print a message and exit
	if (info.Mode() & os.ModeCharDevice) != 0 {
		return nil, fmt.Errorf("no stdin provided, please pipe input or redirect a file")
	}

	// Read all input from stdin
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, fmt.Errorf("error reading from stdin: %w", err)
	}

	// Check if the input is empty
	if len(data) == 0 {
		return nil, fmt.Errorf("stdin is empty")
	}

	return data, nil
}

func ReadFromFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading from file %s, %w", path, err)
	}
	return data, nil
}
