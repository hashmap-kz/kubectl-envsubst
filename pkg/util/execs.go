package util

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type ExecCmdInternalResult struct {
	StdoutContent string
	StderrContent string
}

func ExecCmd(name string, arg ...string) (ExecCmdInternalResult, error) {
	cmd := exec.Command(name, arg...)

	log.Println(name + " " + strings.Join(arg, " "))

	// 1) tee

	//var stdoutBuf, stderrBuf bytes.Buffer
	//cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	//cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)
	//
	//err := cmd.Run()
	//return stdoutBuf.String(), stderrBuf.String(), err

	// 2) silent

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	err := cmd.Run()
	return ExecCmdInternalResult{
		StdoutContent: stdoutBuf.String(),
		StderrContent: stderrBuf.String(),
	}, err
}

func (e ExecCmdInternalResult) CombinedOutput() string {
	return fmt.Sprintf("%s\n%s", e.StdoutContent, e.StderrContent)
}

func ExecWithStdin(name string, stdinContent []byte, arg ...string) (ExecCmdInternalResult, error) {
	// Define the command to execute
	cmd := exec.Command(name, arg...)

	// Create a pipe for stdin
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return ExecCmdInternalResult{}, err
	}

	// Buffers to capture stdout and stderr
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	// Start the command
	if err := cmd.Start(); err != nil {
		return ExecCmdInternalResult{}, err
	}

	// Write to stdin in a separate goroutine
	go func() {
		defer stdin.Close()
		stdin.Write(stdinContent)
	}()

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		return ExecCmdInternalResult{}, err
	}

	return ExecCmdInternalResult{
		StdoutContent: stdoutBuf.String(),
		StderrContent: stderrBuf.String(),
	}, err
}
