package cmd

import (
	"fmt"
	"os"
	"strings"
)

const (
	envsubstAllowedVarsEnv     = "ENVSUBST_ALLOWED_VARS"
	envsubstAllowedPrefixesEnv = "ENVSUBST_ALLOWED_PREFIXES"
)

type CmdArgsRawRecognized struct {
	Filenames             []string
	EnvsubstAllowedVars   []string
	EnvsubstAllowedPrefix []string
	Recursive             bool
	Help                  bool
	Others                []string
	HasStdin              bool
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

func ParseArgs() (CmdArgsRawRecognized, error) {
	args := os.Args[1:] // Skip the program name
	var result CmdArgsRawRecognized

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {

		// --filename=pod.yaml
		case strings.HasPrefix(arg, "--filename="):
			filenameGiven := strings.TrimPrefix(arg, "--filename=")
			if filenameGiven == "" {
				return result, fmt.Errorf("missing value for flag %s", arg)
			}
			err := handleStdin(filenameGiven, &result)
			if err != nil {
				return result, err
			}
			if filenameGiven != "-" {
				result.Filenames = append(result.Filenames, filenameGiven)
			}

			// -f=pod.yaml
		case strings.HasPrefix(arg, "-f="):
			filenameGiven := strings.TrimPrefix(arg, "-f=")
			if filenameGiven == "" {
				return result, fmt.Errorf("missing value for flag %s", arg)
			}
			err := handleStdin(filenameGiven, &result)
			if err != nil {
				return result, err
			}
			if filenameGiven != "-" {
				result.Filenames = append(result.Filenames, filenameGiven)
			}
			// --filename pod.yaml -f pod.yaml
		case arg == "--filename" || arg == "-f":
			if i+1 < len(args) {
				filenameGiven := args[i+1]
				if filenameGiven == "" {
					return result, fmt.Errorf("missing value for flag %s", arg)
				}
				err := handleStdin(filenameGiven, &result)
				if err != nil {
					return result, err
				}
				if filenameGiven != "-" {
					result.Filenames = append(result.Filenames, filenameGiven)
				}
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

	// trying to get allowed-vars config from envs
	if len(result.EnvsubstAllowedVars) == 0 {
		if value, exists := os.LookupEnv(envsubstAllowedVarsEnv); exists {
			split := strings.Split(value, ",")
			if allEmpty(split) {
				return result, fmt.Errorf("missing value for env: %s", envsubstAllowedVarsEnv)
			}
			result.EnvsubstAllowedVars = split
		}
	}

	// trying to get allowed-prefixes from envs
	if len(result.EnvsubstAllowedPrefix) == 0 {
		if value, exists := os.LookupEnv(envsubstAllowedPrefixesEnv); exists {
			split := strings.Split(value, ",")
			if allEmpty(split) {
				return result, fmt.Errorf("missing value for env: %s", envsubstAllowedPrefixesEnv)
			}
			result.EnvsubstAllowedPrefix = split
		}
	}

	return result, nil
}

func handleStdin(filenameGiven string, result *CmdArgsRawRecognized) error {
	if filenameGiven == "-" {
		if result.HasStdin {
			return fmt.Errorf("multiple redirection to stdin detected")
		}
		result.HasStdin = true
	}
	return nil
}
