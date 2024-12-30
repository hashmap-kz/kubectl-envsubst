package util

import (
	"fmt"
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

func resolveFilenames(path string, recursive bool) ([]string, error) {
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
					results = append(results, p)
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
			results = append(results, path)
		}
	}

	// Ensure consistent order
	sort.Strings(results)
	return results, nil
}

// resolveFilenames resolves files, directories, and glob patterns
func resolveFilenames2(inputPath string, recursive bool) ([]string, error) {
	results := []string{}

	matches, err := filepath.Glob(inputPath)
	if err != nil {
		return nil, err
	}

	for _, path := range matches {
		// Check if the path is a directory or file
		_, err := os.Stat(path)
		if err != nil {
			return nil, err
		}
		if os.IsNotExist(err) {
			return nil, err
		}

		visitors, err := ExpandPathsToFileVisitors(path, recursive, FileExtensions)
		if err != nil {
			return nil, err
		}

		results = append(results, visitors...)
	}

	return results, nil
}

func ExpandPathsToFileVisitors(paths string, recursive bool, extensions []string) ([]string, error) {
	visitors := []string{}

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

		visitors = append(visitors, path)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return visitors, nil
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
