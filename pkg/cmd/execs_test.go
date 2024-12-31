package cmd

import "testing"

func TestExecWithStdin(t *testing.T) {
	input := "hello stdin"
	result, err := ExecWithStdin("cat", []byte(input))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.StdoutContent != input {
		t.Errorf("expected stdout to be '%v', got %v", input, result.StdoutContent)
	}

	if result.StderrContent != "" {
		t.Errorf("expected stderr to be empty, got %v", result.StderrContent)
	}
}
