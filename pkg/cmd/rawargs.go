package cmd

import (
	"fmt"
	"os"
	"strings"
)

type CmdArgsRawRecognized struct {
	Filenames             []string
	EnvsubstAllowedVars   []string
	EnvsubstAllowedPrefix []string
	Strict                bool
	Recursive             bool
	Help                  bool
	Others                []string
}

func allEmpty(where []string) bool {
	if len(where) == 0 {
		return true
	}
	for _, s := range where {
		if strings.TrimSpace(s) != "" {
			return false
		}
	}
	return true
}

func parseArgs() (CmdArgsRawRecognized, error) {
	args := os.Args[1:] // Skip the program name
	var result CmdArgsRawRecognized

	// by default working in strict mode
	// may be turned off with --envsubst-no-strict flag
	result.Strict = true

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {

		// --filename=pod.yaml
		case strings.HasPrefix(arg, "--filename="):
			filenameGiven := strings.TrimPrefix(arg, "--filename=")
			if filenameGiven == "" {
				return result, fmt.Errorf("missing value for flag %s", arg)
			}
			result.Filenames = append(result.Filenames, filenameGiven)

			// -f=pod.yaml
		case strings.HasPrefix(arg, "-f="):
			filenameGiven := strings.TrimPrefix(arg, "-f=")
			if filenameGiven == "" {
				return result, fmt.Errorf("missing value for flag %s", arg)
			}
			result.Filenames = append(result.Filenames, filenameGiven)

			// --filename pod.yaml -f pod.yaml
		case arg == "--filename" || arg == "-f":
			if i+1 < len(args) {
				filenameGiven := args[i+1]
				if filenameGiven == "" {
					return result, fmt.Errorf("missing value for flag %s", arg)
				}
				result.Filenames = append(result.Filenames, filenameGiven)
				i++ // Skip the next argument since it's the value
			} else {
				return result, fmt.Errorf("missing value for flag %s", arg)
			}

			// --envsubst-allowed-vars=HOME,USER
		case strings.HasPrefix(arg, "--envsubst-allowed-vars="):
			split := strings.Split(strings.TrimPrefix(arg, "--envsubst-allowed-vars="), ",")
			if allEmpty(split) {
				return result, fmt.Errorf("missing value for flag %s", arg)
			}
			result.EnvsubstAllowedVars = append(result.EnvsubstAllowedVars, split...)

			// --envsubst-allowed-vars HOME,USER
		case arg == "--envsubst-allowed-vars":
			if i+1 < len(args) {
				split := strings.Split(args[i+1], ",")
				if allEmpty(split) {
					return result, fmt.Errorf("missing value for flag %s", arg)
				}
				result.EnvsubstAllowedVars = append(result.EnvsubstAllowedVars, split...)
				i++ // Skip the next argument since it's the value
			} else {
				return result, fmt.Errorf("missing value for flag %s", arg)
			}

			// --envsubst-allowed-prefixes=CI_,APP_
		case strings.HasPrefix(arg, "--envsubst-allowed-prefixes="):
			split := strings.Split(strings.TrimPrefix(arg, "--envsubst-allowed-prefixes="), ",")
			if allEmpty(split) {
				return result, fmt.Errorf("missing value for flag %s", arg)
			}
			result.EnvsubstAllowedPrefix = append(result.EnvsubstAllowedPrefix, split...)

			// --envsubst-allowed-prefixes CI_,APP_
		case arg == "--envsubst-allowed-prefixes":
			if i+1 < len(args) {
				split := strings.Split(args[i+1], ",")
				if allEmpty(split) {
					return result, fmt.Errorf("missing value for flag %s", arg)
				}
				result.EnvsubstAllowedPrefix = append(result.EnvsubstAllowedPrefix, split...)
				i++ // Skip the next argument since it's the value
			} else {
				return result, fmt.Errorf("missing value for flag %s", arg)
			}

		// --envsubst-no-strict
		case arg == "--envsubst-no-strict":
			result.Strict = false

		// --recursive, -R
		case arg == "--recursive" || arg == "-R":
			result.Recursive = true

			// --help, -h
		case arg == "--help" || arg == "-h":
			result.Help = true

		default:
			result.Others = append(result.Others, arg)
		}
	}

	return result, nil
}

func isStrict(mode string) (bool, error) {
	if strings.ToLower(mode) == "not-strict" {
		return false, nil
	}
	if strings.ToLower(mode) == "strict" {
		return true, nil
	}
	return false, fmt.Errorf("incorrect mode: %s, expected one of: strict/not-strict", mode)
}
