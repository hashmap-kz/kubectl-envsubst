package main

import (
	"fmt"
	"github.com/hashmap-kz/kubectl-envsubst/pkg/cmd"
	"os"
)

func main() {

	err := cmd.RunApp()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s", err.Error())
		os.Exit(1)
	}
}
