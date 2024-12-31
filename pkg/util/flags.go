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
	Recursive bool
	Others    []string
}

func ParseCmdFlags(args []string) (*CmdFlagsProxy, error) {
	res := &CmdFlagsProxy{
		Filenames: []string{},
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

		if args[i] == "-f" || args[i] == "--filename" {
			if i+1 >= len(args) {
				return nil, fmt.Errorf("flag --filename requires a value")
			}

			filenameGiven := args[i+1]
			err := resolveSingle(filenameGiven, res)
			if err != nil {
				return nil, err
			}
			i++
			continue
		}

		if strings.HasPrefix(args[i], "-f=") || strings.HasPrefix(args[i], "--filename=") {
			filenameGiven := strings.SplitN(args[i], "=", 2)[1]
			err := resolveSingle(filenameGiven, res)
			if err != nil {
				return nil, err
			}
			continue
		}

		// already handled beforehand, needs to be skipped
		if args[i] == "--recursive" || args[i] == "-R" {
			res.Recursive = true
			continue
		}

		// default
		res.Others = append(res.Others, args[i])
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
