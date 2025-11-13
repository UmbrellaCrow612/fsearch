package runner

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"

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

	var modifiedBefore, modifiedAfter *time.Time

	// --- Parse date filters ---
	if argsMap.ModifiedBefore != "" && argsMap.ModifiedBefore != "Empty" {
		if t, err := time.Parse("2006-01-02", argsMap.ModifiedBefore); err == nil {
			modifiedBefore = &t
		}
	}
	if argsMap.ModifiedAfter != "" && argsMap.ModifiedAfter != "Empty" {
		if t, err := time.Parse("2006-01-02", argsMap.ModifiedAfter); err == nil {
			modifiedAfter = &t
		}
	}

	// --- Convert size filters to bytes ---
	var sizeMultiplier int64 = 1
	switch strings.ToUpper(argsMap.SizeType) {
	case "KB":
		sizeMultiplier = 1024
	case "MB":
		sizeMultiplier = 1024 * 1024
	case "GB":
		sizeMultiplier = 1024 * 1024 * 1024
	case "B", "":
		sizeMultiplier = 1
	default:
		sizeMultiplier = 1
	}

	minSizeBytes := argsMap.MinSize * sizeMultiplier
	maxSizeBytes := argsMap.MaxSize * sizeMultiplier

	// --- Normalize extension filters ---
	normalizeExt := func(exts []string) []string {
		norm := make([]string, 0, len(exts))
		for _, e := range exts {
			e = strings.ToLower(strings.TrimSpace(e))
			if e == "" {
				continue
			}
			if !strings.HasPrefix(e, ".") {
				e = "." + e
			}
			norm = append(norm, e)
		}
		return norm
	}

	includeExts := normalizeExt(argsMap.Ext)
	excludeExts := normalizeExt(argsMap.ExcludeExt)

	hasIncludeExts := len(includeExts) > 0
	hasExcludeExts := len(excludeExts) > 0

	shouldIncludeExt := func(ext string) bool {
		ext = strings.ToLower(ext)
		if hasExcludeExts {
			if slices.Contains(excludeExts, ext) {
				return false
			}
		}
		if hasIncludeExts {
			for _, inc := range includeExts {
				if ext == inc {
					return true
				}
			}
			return false 
		}
		return true 
	}

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
				continue
			}

			info, err := entry.Info()
			if err != nil {
				errMu.Lock()
				errs = append(errs, fmt.Errorf("error getting info for %s: %w", fullPath, err))
				errMu.Unlock()
				continue
			}

			modTime := info.ModTime()
			size := info.Size()
			ext := strings.ToLower(filepath.Ext(entry.Name()))

			// --- Extension filters ---
			if !shouldIncludeExt(ext) {
				continue
			}

			// --- Date filters ---
			if modifiedBefore != nil && modTime.After(*modifiedBefore) {
				continue
			}
			if modifiedAfter != nil && modTime.Before(*modifiedAfter) {
				continue
			}

			// --- Size filters ---
			if argsMap.MinSize > 0 && size < minSizeBytes {
				continue
			}
			if argsMap.MaxSize > 0 && size > maxSizeBytes {
				continue
			}

			mu.Lock()
			files = append(files, fileEntry{Path: fullPath, Entry: entry})
			mu.Unlock()
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
