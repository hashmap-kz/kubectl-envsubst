package cmd

import (
	"testing"
)

func TestExecWithStdin(t *testing.T) {
	t.Run("Successful Execution", func(t *testing.T) {
		// Command and input content
		command := "cat" // `cat` reads from stdin and echoes to stdout
		input := []byte("Hello, World!")

		// Execute the command
		result, err := ExecWithStdin(command, input)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		// Validate stdout
		if result.StdoutContent != string(input) {
			t.Errorf("Expected stdout: %q, got: %q", string(input), result.StdoutContent)
		}

		// Validate stderr
		if result.StderrContent != "" {
			t.Errorf("Expected no stderr, got: %q", result.StderrContent)
		}
	})

	t.Run("Command Error", func(t *testing.T) {
		// Invalid command to simulate error
		command := "nonexistent_command"
		input := []byte("This won't be used")

		// Execute the command
		result, err := ExecWithStdin(command, input)

		// Validate error
		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		// Validate stdout and stderr are empty
		if result.StdoutContent != "" {
			t.Errorf("Expected empty stdout, got: %q", result.StdoutContent)
		}
		if result.StderrContent == "" {
			t.Errorf("Expected stderr to contain error message, got empty")
		}
	})
}
