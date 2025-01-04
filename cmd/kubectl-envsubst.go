package main

import (
	"fmt"
	"github.com/hashmap-kz/kubectl-envsubst/pkg/cmd"
	"os"
	"os/exec"
	"strings"
)

func main() {

	err := runApp()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s", err.Error())
		os.Exit(1)
	}
}

// runApp executes the plugin, with logic divided into smaller, testable components
func runApp() error {

	// parse all passed cmd arguments without any modification
	flags, err := cmd.ParseArgs()
	if err != nil {
		return err
	}

	// either help message, either 'apply' was not provided
	if flags.Help || len(flags.Others) == 0 {
		fmt.Println(cmd.UsageMessage())
		return nil
	}

	// support apply operation only
	if flags.Others[0] != "apply" {
		fmt.Println(cmd.UsageMessage())
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

	// prepare files for apply all of them as a single stream
	joinedFilesData, err := cmd.JoinFiles(files, flags.HasStdin)
	if err != nil {
		return err
	}

	// substitute the whole stream of joined files at once
	envSubst := cmd.NewEnvsubst(flags.EnvsubstAllowedVars, flags.EnvsubstAllowedPrefix, true)
	streams, err := envSubst.SubstituteEnvs(string(joinedFilesData))
	if err != nil {
		return err
	}

	// prepare kubectl args
	args := []string{}
	args = append(args, flags.Others...)
	args = append(args, "-f")
	args = append(args, "-")

	// pass stream of files to stdin
	execCmd, err := cmd.ExecWithStdin(kubectl, []byte(streams), args...)
	if err != nil {
		fmt.Println(strings.TrimSpace(execCmd.StderrContent))
		return err
	}

	fmt.Println(strings.TrimSpace(execCmd.StdoutContent))
	return nil

}
