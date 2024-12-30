package util

import (
	"errors"
	"fmt"
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
				files, err := resolveFilenames([]string{args[i+1]}, res.Recursive)
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

func resolveFilenames(inputPaths []string, recursive bool) ([]string, error) {
	results := []string{}
	for _, s := range inputPaths {
		switch {

		case strings.Index(s, "http://") == 0 || strings.Index(s, "https://") == 0:
			url, err := url.Parse(s)
			if err != nil {
				return nil, err
			}
			results = append(results, url.String())
		default:
			matches, err := expandIfFilePattern(s)
			if err != nil {
				return nil, err
			}

			builderPaths, err := iterateOverMatches(recursive, matches...)
			if err != nil {
				return nil, err
			}
			results = append(results, builderPaths...)
		}
	}

	return results, nil
}

func iterateOverMatches(recursive bool, paths ...string) ([]string, error) {
	results := []string{}

	for _, p := range paths {
		_, err := os.Stat(p)
		if os.IsNotExist(err) {
			return nil, err
		}
		if err != nil {
			return nil, err
		}

		expandedPaths, err := expandPaths(p, recursive, FileExtensions)
		if err != nil {
			return nil, err
		}
		results = append(results, expandedPaths...)
	}

	return results, nil
}

func expandPaths(paths string, recursive bool, extensions []string) ([]string, error) {
	results := []string{}

	err := filepath.Walk(paths, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fi.IsDir() {
			if path != paths && !recursive {
				return filepath.SkipDir
			}
			return nil
		}
		// Don't check extension if the filepath was passed explicitly
		if path != paths && ignoreFile(path, extensions) {
			return nil
		}

		results = append(results, path)
		return nil
	})

	if err != nil {
		return nil, err
	}
	return results, nil
}

func expandIfFilePattern(pattern string) ([]string, error) {
	if _, err := os.Stat(pattern); os.IsNotExist(err) {
		matches, err := filepath.Glob(pattern)
		if err == nil && len(matches) == 0 {
			return nil, fmt.Errorf("path not exist: %s", pattern)
		}
		if errors.Is(err, filepath.ErrBadPattern) {
			return nil, fmt.Errorf("pattern %q is not valid: %v", pattern, err)
		}
		return matches, err
	}
	return []string{pattern}, nil
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

func resolveFilenames2(path string, recursive bool) ([]string, error) {
	var results []string

	// Handle glob patterns
	if strings.Contains(path, "*") {
		matches, err := filepath.Glob(path)
		if err != nil {
			return nil, err
		}
		results = append(results, matches...)
	} else {
		// Check if the path is a directory or file
		info, err := os.Stat(path)
		if err != nil {
			return nil, err
		}

		if info.IsDir() {
			// List files in the directory
			err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					results = append(results, filepath.Clean(p))
				}
				if !recursive && info.IsDir() && p != path {
					return filepath.SkipDir
				}
				return nil
			})
			if err != nil {
				return nil, err
			}
		} else {
			results = append(results, filepath.Clean(path))
		}
	}

	// Ensure consistent order
	sort.Strings(results)
	return results, nil
}
