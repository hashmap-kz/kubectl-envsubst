package util

import (
	"fmt"
)

// CmdFlagsProxy holds interested for plugin args
type CmdFlagsProxy struct {
	Filenames []string
	Namespace string
	Others    []string
}

func ParseCmdFlags(args []string) (*CmdFlagsProxy, error) {
	res := &CmdFlagsProxy{
		Filenames: []string{},
		Namespace: "default",
		Others:    []string{},
	}

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--filename", "-f":
			if i+1 < len(args) {
				res.Filenames = append(res.Filenames, args[i+1])
				i++
			} else {
				return nil, fmt.Errorf("flag --filename requires a value")
			}

		case "--namespace", "-n":
			if i+1 < len(args) {
				res.Namespace = args[i+1]
				i++
			} else {
				return nil, fmt.Errorf("flag --namespace requires a value")
			}

		default:
			res.Others = append(res.Others, args[i])
		}

	}

	return res, nil
}
