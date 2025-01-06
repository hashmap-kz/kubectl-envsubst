package main

import (
	"fmt"
	"github.com/hashmap-kz/kubectl-envsubst/cmd/app"
	"os"
)

func main() {

	err := app.RunApp()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s", err.Error())
		os.Exit(1)
	}
}
