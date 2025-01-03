package cmd

import (
	"bytes"
	"os/exec"
)

type ExecCmdInternalResult struct {
	StdoutContent string
	StderrContent string
}

func ExecWithStdin(name string, stdinContent []byte, arg ...string) (ExecCmdInternalResult, error) {
	// Define the command to execute
	cmd := exec.Command(name, arg...)

	// Buffers to capture stdout and stderr
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	// Create a pipe for stdin
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return ExecCmdInternalResult{
			StdoutContent: stdoutBuf.String(),
			StderrContent: getErrorDesc(err, stderrBuf),
		}, err
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return ExecCmdInternalResult{
			StdoutContent: stdoutBuf.String(),
			StderrContent: getErrorDesc(err, stderrBuf),
		}, err
	}

	// Write to stdin in a separate goroutine
	go func() {
		defer stdin.Close()
		stdin.Write(stdinContent)
	}()

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		return ExecCmdInternalResult{
			StdoutContent: stdoutBuf.String(),
			StderrContent: getErrorDesc(err, stderrBuf),
		}, err
	}

	return ExecCmdInternalResult{
		StdoutContent: stdoutBuf.String(),
		StderrContent: stderrBuf.String(),
	}, err
}

func getErrorDesc(err error, stderrBuf bytes.Buffer) string {
	errorMessage := err.Error()
	if stderrBuf.String() != "" {
		errorMessage = stderrBuf.String()
	}
	return errorMessage
}
