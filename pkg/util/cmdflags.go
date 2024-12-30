package util

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
	Filenames []string
	Namespace string
	Others    []string
	Recursive bool
}

func ParseCmdFlags(args []string) (*CmdFlagsProxy, error) {
	res := &CmdFlagsProxy{
		Filenames: []string{},
		Namespace: "default",
		Others:    []string{},
	}

	// handle pre-flight args that needed beforehand
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--recursive", "-R":
			res.Recursive = true
		}
	}

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--filename", "-f":
			if i+1 < len(args) {
				files, err := resolveFilenames(args[i+1], res.Recursive)
				if err != nil {
					return nil, fmt.Errorf("error resolving filenames: %w", err)
				}
				res.Filenames = append(res.Filenames, files...)
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

			// already handled, but needs to be skipped
		case "--recursive", "-R":
			res.Recursive = true

		default:
			res.Others = append(res.Others, args[i])
		}

	}

	return res, nil
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
