package cmd

import (
	"bytes"
	"io"
	"os"
)

func JoinFiles(flags *CmdFlagsProxy) ([]byte, error) {
	buf := bytes.Buffer{}

	totalFiles := len(flags.Filenames)
	if flags.HasStdin {
		totalFiles += 1
	}
	needSeparator := totalFiles > 1

	// process STDIN
	if flags.HasStdin {
		stdin, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, err
		}

		substituted, err := substBuf(flags, stdin)
		if err != nil {
			return nil, err
		}

		buf.WriteString(substituted)
		if needSeparator {
			buf.WriteString("\n---\n")
		}
	}

	// process files
	for _, f := range flags.Filenames {

		// get file data
		file, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}

		substituted, err := substBuf(flags, file)
		if err != nil {
			return nil, err
		}

		buf.WriteString(substituted)
		if needSeparator {
			buf.WriteString("\n---\n")
		}
	}

	return buf.Bytes(), nil
}

func substBuf(flags *CmdFlagsProxy, data []byte) (string, error) {
	// substitute environment variables
	// strict mode is always ON
	envSubst := NewEnvsubst(flags.EnvsubstAllowedVars, flags.EnvsubstAllowedPrefix, true)
	substituted, err := envSubst.SubstituteEnvs(string(data))
	if err != nil {
		return "", err
	}
	return substituted, nil
}
