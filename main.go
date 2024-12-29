package main

import (
	"flag"
	"fmt"
	"github.com/hashmap-kz/kubectl-envsubst/pkg/app"
	"github.com/hashmap-kz/kubectl-envsubst/pkg/util"
	"log"
	"os"
	"strings"
)

func main() {
	// Define required flags
	allowedFlag := flag.String("a", "", "Comma-separated list of allowed variable names")
	inputFlag := flag.String("f", "", "Input file")
	flag.Parse()

	// Helper function to check if a flag is set
	isFlagSet := func(name string) bool {
		return flag.Lookup(name).Value.String() != flag.Lookup(name).DefValue
	}

	var data []byte
	var err error

	if !isFlagSet("f") {
		fmt.Fprintln(os.Stderr, "Input is not specified")
		os.Exit(1)
	}

	if *inputFlag == "-" || *inputFlag == "" {
		data, err = util.ReadFromStdin()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		data, err = util.ReadFromFile(*inputFlag)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Get the list of allowed variable names
	allowedList := strings.Split(*allowedFlag, ",")

	// Substitute and print result to stdout
	envs := app.SubstituteEnvs(string(data), allowedList)
	fmt.Println(envs)

}
