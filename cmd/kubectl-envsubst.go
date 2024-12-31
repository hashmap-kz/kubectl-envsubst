package main

import (
	"fmt"
	"github.com/hashmap-kz/kubectl-envsubst/pkg/cmd"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {

	flags, err := cmd.ParseCmdFlags()
	if err != nil {
		log.Fatal(err)
	}

	// either help message, either 'apply' was not provided
	if flags.Help || len(flags.Others) == 0 {
		fmt.Println(cmd.UsageMessage())
		os.Exit(0)
	}

	// support apply operation only
	if flags.Others[0] != "apply" {
		fmt.Println(cmd.UsageMessage())
		os.Exit(0)
	}

	kubectl, err := exec.LookPath("kubectl")
	if err != nil {
		log.Fatal(err)
	}

	// prepare files to apply as a single stream
	streams, err := cmd.JoinFiles(flags)
	if err != nil {
		log.Fatal(err)
	}

	// prepare kubectl args
	args := []string{}
	args = append(args, flags.Others...)
	args = append(args, "-f")
	args = append(args, "-")

	// pass stream of files to stdin
	cmd, err := cmd.ExecWithStdin(kubectl, streams, args...)
	if err != nil {
		fmt.Println(strings.TrimSpace(cmd.StderrContent))
		log.Fatal(err)
	}

	fmt.Println(strings.TrimSpace(cmd.StdoutContent))

}
