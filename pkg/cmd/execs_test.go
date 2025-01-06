package cmd

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
)

// Cross-platform helper to get a command
func getCommand(command string, args []string) (c string, a []string) {
	if runtime.GOOS == "windows" {
		// Adjust commands for Windows
		switch command {
		case "cat":
			// Use "more" to simulate "cat" on Windows for stdin
			return "more", args
		case "sh":
			return "cmd", append([]string{"/C"}, args...)
		case "ls":
			return "cmd", append([]string{"/C", "dir"}, args...)
		default:
			return command, args
		}
	}
	return command, args // Default for Linux/macOS
}

func TestExecWithStdin(t *testing.T) {
	t.Run("Successful Execution", func(t *testing.T) {
		// Command and input content
		command, args := getCommand("cat", nil)
		input := []byte("Hello, World!")

		// Execute the command
		result, err := ExecWithStdin(command, input, args...)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		// Validate stdout
		if strings.TrimSpace(result.StdoutContent) != strings.TrimSpace(string(input)) {
			t.Errorf("Expected stdout: %q, got: %q", string(input), result.StdoutContent)
		}

		// Validate stderr
		if result.StderrContent != "" {
			t.Errorf("Expected no stderr, got: %q", result.StderrContent)
		}
	})

	t.Run("Command Error", func(t *testing.T) {
		// Invalid command to simulate error
		command, args := getCommand("nonexistent_command", nil)
		input := []byte("This won't be used")

		// Execute the command
		result, err := ExecWithStdin(command, input, args...)

		// Validate error
		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		// Validate stdout and stderr
		if result.StdoutContent != "" {
			t.Errorf("Expected empty stdout, got: %q", result.StdoutContent)
		}
		if result.StderrContent == "" {
			t.Errorf("Expected stderr to contain error message, got empty")
		}
	})

	t.Run("Command Produces Stderr", func(t *testing.T) {
		// Command that writes to stderr
		command, args := getCommand("sh", []string{"-c", "echo 'error message' 1>&2"})
		if runtime.GOOS == "windows" {
			command, args = getCommand("cmd", []string{"/C", "echo error message 1>&2"})
		}
		input := []byte("")

		// Execute the command
		result, err := ExecWithStdin(command, input, args...)
		// Validate error (should not fail because the command exits with 0)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		// Validate stdout
		if result.StdoutContent != "" {
			t.Errorf("Expected empty stdout, got: %q", result.StdoutContent)
		}

		// Validate stderr
		expectedStderr := "error message"
		if strings.TrimSpace(result.StderrContent) != strings.TrimSpace(expectedStderr) {
			t.Errorf("Expected stderr: %q, got: %q", expectedStderr, result.StderrContent)
		}
	})

	t.Run("Command With Invalid Arguments", func(t *testing.T) {
		// Command with invalid arguments
		command, args := getCommand("ls", []string{"--invalid-option"})
		input := []byte("")

		// Execute the command
		result, err := ExecWithStdin(command, input, args...)

		// Validate error
		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		// Validate stdout and stderr
		if result.StdoutContent != "" {
			t.Errorf("Expected empty stdout, got: %q", result.StdoutContent)
		}
		if result.StderrContent == "" {
			t.Errorf("Expected stderr to contain error message, got empty")
		}
	})
}

func TestExecWithLargeInput(t *testing.T) {
	// Generate a large input string
	var largeInput strings.Builder
	for i := 0; i < 1000000; i++ { // 1 million lines
		largeInput.WriteString("Line " + fmt.Sprint(i) + "\n")
	}

	// Execute the command with the large input
	result, err := ExecWithStdin("cat", []byte(largeInput.String()))
	// Validate that the command executed successfully
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Validate that the output matches the input
	if result.StdoutContent != largeInput.String() {
		t.Errorf("Output does not match input\n")
	}

	// Validate stderr is empty
	if result.StderrContent != "" {
		t.Errorf("Expected no stderr, got: %s", result.StderrContent)
	}
}
