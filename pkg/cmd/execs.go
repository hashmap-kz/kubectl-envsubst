package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
)

type ExecCmdInternalResult struct {
	StdoutContent string
	StderrContent string
}

func ExecWithStdin(name string, stdinContent []byte, arg ...string) (ExecCmdInternalResult, error) {
	cmd := exec.Command(name, arg...)

	// Buffers to capture stdout and stderr
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	// Create a pipe for stdin
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return resultFromError(err, stderrBuf)
	}

	// Create a channel to capture errors from the goroutine
	writeErrChan := make(chan error, 1)

	// Write to stdin in a separate goroutine
	go func() {
		defer close(writeErrChan)
		_, writeErr := stdin.Write(stdinContent)
		if writeErr != nil {
			writeErrChan <- writeErr
			return
		}
		writeErrChan <- stdin.Close()
	}()

	// Start the command
	if err := cmd.Start(); err != nil {
		return resultFromError(err, stderrBuf)
	}

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		return resultFromError(err, stderrBuf)
	}

	// Check if the write to stdin failed
	if writeErr := <-writeErrChan; writeErr != nil {
		return resultFromError(writeErr, stderrBuf)
	}

	return ExecCmdInternalResult{
		StdoutContent: stdoutBuf.String(),
		StderrContent: stderrBuf.String(),
	}, nil
}

func getErrorDesc(err error, stderrBuf bytes.Buffer) string {
	if stderrBuf.Len() > 0 {
		return stderrBuf.String()
	}
	return err.Error()
}

func resultFromError(err error, stderrBuf bytes.Buffer) (ExecCmdInternalResult, error) {
	return ExecCmdInternalResult{
		StderrContent: getErrorDesc(err, stderrBuf),
	}, fmt.Errorf("execution failed: %w", err)
}
