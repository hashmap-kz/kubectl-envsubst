package main

import (
	"flag"
	"fmt"
	"github.com/hashmap-kz/kubectl-envsubst/pkg/app"
	"io"
	"os"
	"strings"
)

func main() {
	// Define the --allowed flag
	allowedFlag := flag.String("a", "", "Comma-separated list of allowed variable names")
	flag.Parse()

	// Get the list of allowed variable names
	allowedList := strings.Split(*allowedFlag, ",")

	// Example input string (from stdin or other sources)
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading from stdin:", err)
		os.Exit(1)
	}

	envs := app.SubstituteEnvs(string(data), allowedList)
	fmt.Println(envs)
}
