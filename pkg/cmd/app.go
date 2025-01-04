package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// RunApp executes the plugin, with logic divided into smaller, testable components
func RunApp() error {

	// parse all passed cmd arguments without any modification
	flags, err := ParseArgs()
	if err != nil {
		return err
	}

	// either help message, either 'apply' was not provided
	if flags.Help || len(flags.Others) == 0 {
		fmt.Println(UsageMessage())
		return nil
	}

	// support apply operation only
	if flags.Others[0] != "apply" {
		fmt.Println(UsageMessage())
		return nil
	}

	// it checks that executable exists
	kubectl, err := exec.LookPath("kubectl")
	if err != nil {
		return err
	}

	// resolve all filenames: expand all glob-patterns, list directories, etc...
	files, err := ResolveAllFiles(flags.Filenames, flags.Recursive)
	if err != nil {
		return err
	}

	// apply STDIN (if any)
	if flags.HasStdin {
		err := applyStdin(flags, kubectl)
		if err != nil {
			return err
		}
	}

	// apply passed files
	for _, filename := range files {
		err := applyOneFile(flags, kubectl, filename)
		if err != nil {
			return err
		}
	}

	return nil
}

// applyStdin substitutes content, passed to stdin `kubectl apply -f -`
func applyStdin(flags CmdArgsRawRecognized, kubectl string) error {
	stdin, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	// substitute the whole stream of joined files at once
	substitutedBuffer, err := substituteContent(flags, stdin)
	if err != nil {
		return err
	}

	return execKubectl(flags, kubectl, substitutedBuffer)
}

// applyOneFile read file (url, local-path), substitute its content, apply result
func applyOneFile(flags CmdArgsRawRecognized, kubectl string, filename string) error {

	// recognize file type

	var contentForSubst []byte
	if IsURL(filename) {
		data, err := readRemoteFileContent(filename)
		if err != nil {
			return err
		}
		contentForSubst = data
	} else {
		data, err := os.ReadFile(filename)
		if err != nil {
			return err
		}
		contentForSubst = data
	}

	// substitute the whole stream of joined files at once
	substitutedBuffer, err := substituteContent(flags, contentForSubst)
	if err != nil {
		return err
	}

	return execKubectl(flags, kubectl, substitutedBuffer)
}

// substituteContent runs the subst module for a given content
func substituteContent(flags CmdArgsRawRecognized, contentForSubst []byte) (string, error) {
	envSubst := NewEnvsubst(flags.EnvsubstAllowedVars, flags.EnvsubstAllowedPrefix, true)
	substitutedBuffer, err := envSubst.SubstituteEnvs(string(contentForSubst))
	if err != nil {
		return "", err
	}
	return substitutedBuffer, nil
}

// execKubectl applies a result buffer, bu running `kubectl apply -f -`
func execKubectl(flags CmdArgsRawRecognized, kubectl string, substitutedBuffer string) error {
	// prepare kubectl args
	args := []string{}
	args = append(args, flags.Others...)
	args = append(args, "-f")
	args = append(args, "-")

	// pass stream of files to stdin
	execCmd, err := ExecWithStdin(kubectl, []byte(substitutedBuffer), args...)
	if err != nil {
		fmt.Println(strings.TrimSpace(execCmd.StderrContent))
		return err
	}

	fmt.Println(strings.TrimSpace(execCmd.StdoutContent))
	return nil
}
