package runner

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/UmbrellaCrow612/fsearch/src/args"
	"github.com/UmbrellaCrow612/fsearch/src/out"
	"github.com/UmbrellaCrow612/fsearch/src/utils"
)

// Runs the main loop and logic
func Run(argsMap *args.ArgsMap) {
	var searchTermRegex *regexp.Regexp
	if argsMap.Regex {
		regex, err := utils.CompileRegex(argsMap.Term)
		if err != nil {
			out.ExitError(err.Error())
		}

		searchTermRegex = regex
	} else {
		regex, err := utils.BuildSearchRegex(argsMap.Term, argsMap.Partial, argsMap.IgnoreCase)
		if err != nil {
			out.ExitError(err.Error())
		}

		searchTermRegex = regex
	}

	if argsMap.Type == "file" {

		// read files pass argmap
		_, err := readFilesParallel(argsMap.Path, argsMap, searchTermRegex)
		if err != nil {
			out.ExitError(err.Error())
		}
		// collect based on flags
		// get collection
		// apply addiotnal flags
		// print end
	} else {
		// read folder
		// collect based on flags
		// get collection
		// apply addiotnal flags
		// print end
	}
}

const maxWorkers = 10

type fileEntry struct {
	Path  string
	Entry os.DirEntry
}

// Reads all *files* in a directory tree in parallel
func readFilesParallel(root string, argsMap *args.ArgsMap, searchTermRegex *regexp.Regexp) ([]fileEntry, error) {
	var files []fileEntry
	var mu sync.Mutex
	var wg sync.WaitGroup
	var errs []error
	var errMu sync.Mutex
	sem := make(chan struct{}, maxWorkers)

	var read func(string)
	read = func(path string) {
		sem <- struct{}{}
		defer func() { <-sem }()
		defer wg.Done()

		list, err := os.ReadDir(path)
		if err != nil {
			errMu.Lock()
			errs = append(errs, fmt.Errorf("error reading %s: %w", path, err))
			errMu.Unlock()
			return
		}

		for _, entry := range list {
			fullPath := filepath.Join(path, entry.Name())

			if entry.IsDir() {
				wg.Add(1)
				go read(fullPath)
			} else {
				mu.Lock()
				files = append(files, fileEntry{Path: fullPath, Entry: entry})
				mu.Unlock()
			}
		}
	}

	wg.Add(1)
	go read(root)
	wg.Wait()

	if len(errs) > 0 {
		return files, errors.Join(errs...)
	}
	return files, nil
}

// Reads all *directories* in a directory tree in parallel
func readDirsParallel(root string, argsMap *args.ArgsMap) ([]fileEntry, error) {
	var dirs []fileEntry
	var mu sync.Mutex
	var wg sync.WaitGroup
	var errs []error
	var errMu sync.Mutex
	sem := make(chan struct{}, maxWorkers)

	var read func(string)
	read = func(path string) {
		sem <- struct{}{}
		defer func() { <-sem }()
		defer wg.Done()

		list, err := os.ReadDir(path)
		if err != nil {
			errMu.Lock()
			errs = append(errs, fmt.Errorf("error reading %s: %w", path, err))
			errMu.Unlock()
			return
		}

		for _, entry := range list {
			fullPath := filepath.Join(path, entry.Name())

			if entry.IsDir() {
				mu.Lock()
				dirs = append(dirs, fileEntry{Path: fullPath, Entry: entry})
				mu.Unlock()

				wg.Add(1)
				go read(fullPath)
			}
		}
	}

	wg.Add(1)
	go read(root)
	wg.Wait()

	if len(errs) > 0 {
		return dirs, errors.Join(errs...)
	}
	return dirs, nil
}
