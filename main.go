package main

import (
	"fmt"
	"github.com/hashmap-kz/kubectl-envsubst/pkg/util"
	"log"
	"os"
	"strings"
)

func main() {
	flags, err := util.ParseCmdFlags(os.Args[1:])
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

		// pass to stdin
		cmd, err := util.ExecWithStdin("kubectl", file, args...)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(strings.TrimSpace(cmd.StdoutContent))
	}

}
