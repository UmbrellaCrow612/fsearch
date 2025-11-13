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
	"sync/atomic"
	"time"

	"github.com/UmbrellaCrow612/fsearch/src/args"
	"github.com/UmbrellaCrow612/fsearch/src/out"
	"github.com/UmbrellaCrow612/fsearch/src/shared"
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

	matches, err := readInParallel(argsMap.Path, argsMap, searchTermRegex)
	if err != nil {
		out.ExitError(err.Error())
	}

	if argsMap.Open && len(matches) > 0 {
		err := utils.OpenMatchEntry(matches[0])
		if err != nil {
			out.ExitError(err.Error())
		}

		out.ExitSuccess()
	}
}

const maxWorkers = 10

// Reads in a directory tree in parallel
func readInParallel(root string, argsMap *args.ArgsMap, searchTermRegex *regexp.Regexp) ([]shared.MatchEntry, error) {
	var matchEntrys []shared.MatchEntry
	var mu sync.Mutex
	var wg sync.WaitGroup
	var errs []error
	var errMu sync.Mutex
	sem := make(chan struct{}, maxWorkers)

	var modifiedBefore, modifiedAfter *time.Time

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

	sizeMultiplier := utils.GetSizeMultipler(argsMap)

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
			return slices.Contains(includeExts, ext)
		}
		return true
	}

	maxDepth := argsMap.Depth
	if maxDepth <= 0 {
		maxDepth = -1 // no limit
	}

	rootDepth := len(strings.Split(filepath.Clean(root), string(os.PathSeparator)))

	limit := argsMap.Limit
	if limit <= 0 {
		limit = -1 // no limit
	}
	var reachedLimit atomic.Bool

	var read func(string, int, bool)
	read = func(path string, depth int, dirsOnly bool) {
		if limit > 0 && reachedLimit.Load() {
			return
		}

		sem <- struct{}{}
		defer func() { <-sem }()
		defer wg.Done()

		if maxDepth > 0 && depth-rootDepth > maxDepth {
			return
		}

		list, err := os.ReadDir(path)
		if err != nil {
			errMu.Lock()
			errs = append(errs, fmt.Errorf("error reading %s: %w", path, err))
			errMu.Unlock()
			return
		}

		for _, entry := range list {
			if limit > 0 && reachedLimit.Load() {
				return
			}

			fullPath := filepath.Join(path, entry.Name())

			if entry.IsDir() {
				if !argsMap.Hidden && utils.IsHiddenFolderName(entry.Name()) {
					continue
				}

				if slices.Contains(argsMap.ExcludeDir, entry.Name()) {
					continue
				}

				if dirsOnly {
					if searchTermRegex.MatchString(entry.Name()) {
						mu.Lock()
						if limit > 0 && len(matchEntrys) >= limit {
							mu.Unlock()
							reachedLimit.Store(true)
							return
						}
						matchEntrys = append(matchEntrys, shared.MatchEntry{Path: fullPath, Entry: entry})
						mu.Unlock()
					}

				}

				if maxDepth <= 0 || depth-rootDepth < maxDepth {
					wg.Add(1)
					go read(fullPath, depth+1, dirsOnly)
				}

				continue // Skip file-processing logic
			}

			if dirsOnly {
				continue // If we only want directories, skip all file logic
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

			if !searchTermRegex.MatchString(fullPath) {
				continue
			}

			// --- Add file if limit not reached ---
			mu.Lock()
			if limit > 0 && len(matchEntrys) >= limit {
				mu.Unlock()
				reachedLimit.Store(true)
				return
			}
			matchEntrys = append(matchEntrys, shared.MatchEntry{Path: fullPath, Entry: entry})
			mu.Unlock()
		}
	}

	wg.Add(1)
	go read(root, rootDepth, argsMap.Type == "folder")
	wg.Wait()

	if len(errs) > 0 {
		return matchEntrys, errors.Join(errs...)
	}
	return matchEntrys, nil
}
