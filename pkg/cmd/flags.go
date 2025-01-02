package cmd

import (
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var FileExtensions = []string{".json", ".yaml", ".yml"}

// CmdFlagsProxy holds interested for plugin args
type CmdFlagsProxy struct {
	Filenames             []string
	EnvsubstAllowedVars   []string
	EnvsubstAllowedPrefix []string
	Recursive             bool
	Help                  bool
	Others                []string
}

func ParseCmdFlags() (*CmdFlagsProxy, error) {

	recognized, err := parseArgs()
	if err != nil {
		return nil, err
	}

	res := &CmdFlagsProxy{
		Filenames:             []string{},
		EnvsubstAllowedVars:   recognized.EnvsubstAllowedVars,
		EnvsubstAllowedPrefix: recognized.EnvsubstAllowedPrefix,
		Recursive:             recognized.Recursive,
		Help:                  recognized.Help,
		Others:                recognized.Others,
	}

	for _, f := range recognized.Filenames {
		err := resolveSingle(f, res)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func resolveSingle(filenameGiven string, res *CmdFlagsProxy) error {
	if filenameGiven == "" {
		return fmt.Errorf("flag --filename requires a value")
	}
	files, err := resolveFilenames(filenameGiven, res.Recursive)
	if err != nil {
		return fmt.Errorf("error resolving filenames: %w", err)
	}
	res.Filenames = append(res.Filenames, files...)
	return nil
}

func ignoreFile(path string, extensions []string) bool {
	if len(extensions) == 0 {
		return false
	}
	ext := filepath.Ext(path)
	for _, s := range extensions {
		if s == ext {
			return false
		}
	}
	return true
}

func resolveFilenames(path string, recursive bool) ([]string, error) {
	var results []string

	// Check if the path is a URL
	isURL := func(s string) bool {
		u, err := url.Parse(s)
		return err == nil && u.Scheme != "" && u.Host != ""
	}

	if isURL(path) {
		// Add URL directly to results
		results = append(results, path)
	} else if strings.Contains(path, "*") {
		// Handle glob patterns
		matches, err := filepath.Glob(path)
		if err != nil {
			return nil, fmt.Errorf("error resolving glob pattern: %w", err)
		}
		results = append(results, matches...)
	} else {
		// Check if the path is a directory or file
		info, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("error accessing path: %w", err)
		}

		if info.IsDir() {
			// Walk the directory
			err := filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() && !recursive && p != path {
					return filepath.SkipDir
				}
				if !d.IsDir() {
					if !ignoreFile(filepath.Clean(p), FileExtensions) {
						results = append(results, filepath.Clean(p))
					}
				}
				return nil
			})
			if err != nil {
				return nil, fmt.Errorf("error walking directory: %w", err)
			}
		} else {
			if !ignoreFile(filepath.Clean(path), FileExtensions) {
				results = append(results, filepath.Clean(path))
			}
		}
	}

	// Ensure consistent order
	sort.Strings(results)
	return results, nil
}
