package main

import (
	"fmt"
	"os"

	"github.com/hashmap-kz/kubectl-envsubst/cmd/app"
)

func main() {
	err := app.RunApp()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s", err.Error())
		os.Exit(1)
	}
}
