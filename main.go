package main

import (
	"fmt"
	"github.com/hashmap-kz/kubectl-envsubst/pkg/app"
	"github.com/hashmap-kz/kubectl-envsubst/pkg/util"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {

	flags, err := util.ParseCmdFlags()
	if err != nil {
		log.Fatal(err)
	}

	if flags.Help {
		fmt.Println(util.UsageMessage())
		os.Exit(0)
	}

	kubectl, err := exec.LookPath("kubectl")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range flags.Filenames {
		// prepare kubectl args
		args := []string{}
		args = append(args, flags.Others...)
		args = append(args, "-f")
		args = append(args, "-")

		// get file data
		file, err := os.ReadFile(f)
		if err != nil {
			log.Fatal(err)
		}

		// substitute environment variables
		envSubst := app.NewEnvsubst(flags.EnvsubstAllowedVars, flags.EnvsubstAllowedPrefix, flags.Strict)
		substituted, err := envSubst.SubstituteEnvs(string(file))
		if err != nil {
			log.Fatal(err)
		}

		// pass to stdin
		cmd, err := util.ExecWithStdin(kubectl, []byte(substituted), args...)
		if err != nil {
			fmt.Println(strings.TrimSpace(cmd.StderrContent))
			log.Fatal(err)
		}

		fmt.Println(strings.TrimSpace(cmd.StdoutContent))
	}
}
