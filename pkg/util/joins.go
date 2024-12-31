package util

import (
	"bytes"
	"github.com/hashmap-kz/kubectl-envsubst/pkg/app"
	"os"
)

func JoinFiles(flags *CmdFlagsProxy) ([]byte, error) {
	buf := bytes.Buffer{}

	for _, f := range flags.Filenames {

		// get file data
		file, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}

		// substitute environment variables
		envSubst := app.NewEnvsubst(flags.EnvsubstAllowedVars, flags.EnvsubstAllowedPrefix, flags.Strict)
		substituted, err := envSubst.SubstituteEnvs(string(file))
		if err != nil {
			return nil, err
		}

		buf.WriteString(substituted)
		if len(flags.Filenames) > 1 {
			buf.WriteString("\n---\n")
		}
	}

	return buf.Bytes(), nil
}
