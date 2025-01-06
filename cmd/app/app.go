package app

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/hashmap-kz/kubectl-envsubst/pkg/cmd"
)

// runApp executes the plugin, with logic divided into smaller, testable components
func RunApp() error {
	// parse all passed cmd arguments without any modification
	flags, err := cmd.ParseArgs()
	if err != nil {
		return err
	}

	// either help message, either 'apply' was not provided
	if flags.Help || len(flags.Others) == 0 {
		fmt.Println(cmd.UsageMessage)
		return nil
	}

	// support apply operation only
	if flags.Others[0] != "apply" {
		fmt.Println(cmd.UsageMessage)
		return nil
	}

	// it checks that executable exists
	kubectl, err := exec.LookPath("kubectl")
	if err != nil {
		return err
	}

	// resolve all filenames: expand all glob-patterns, list directories, etc...
	files, err := cmd.ResolveAllFiles(flags.Filenames, flags.Recursive)
	if err != nil {
		return err
	}

	// apply STDIN (if any)
	if flags.HasStdin {
		err := applyStdin(&flags, kubectl)
		if err != nil {
			return err
		}
	}

	// apply passed files
	for _, filename := range files {
		err := applyOneFile(&flags, kubectl, filename)
		if err != nil {
			return err
		}
	}

	return nil
}

// applyStdin substitutes content, passed to stdin `kubectl apply -f -`
func applyStdin(flags *cmd.ArgsRawRecognized, kubectl string) error {
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
func applyOneFile(flags *cmd.ArgsRawRecognized, kubectl, filename string) error {
	// recognize file type

	var contentForSubst []byte
	if cmd.IsURL(filename) {
		data, err := cmd.ReadRemoteFileContent(filename)
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
func substituteContent(flags *cmd.ArgsRawRecognized, contentForSubst []byte) (string, error) {
	envSubst := cmd.NewEnvsubst(flags.EnvsubstAllowedVars, flags.EnvsubstAllowedPrefix, true)
	substitutedBuffer, err := envSubst.SubstituteEnvs(string(contentForSubst))
	if err != nil {
		return "", err
	}
	return substitutedBuffer, nil
}

// execKubectl applies a result buffer, bu running `kubectl apply -f -`
func execKubectl(flags *cmd.ArgsRawRecognized, kubectl, substitutedBuffer string) error {
	// prepare kubectl args
	args := []string{}
	args = append(args, flags.Others...)
	args = append(args, "-f", "-")

	// pass stream of files to stdin
	execCmd, err := cmd.ExecWithStdin(kubectl, []byte(substitutedBuffer), args...)
	if err != nil {
		fmt.Println(strings.TrimSpace(execCmd.StderrContent))
		return err
	}

	fmt.Println(strings.TrimSpace(execCmd.StdoutContent))
	return nil
}
