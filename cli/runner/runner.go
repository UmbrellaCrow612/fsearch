package runner

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/UmbrellaCrow612/fsearch/cli/args"
	"github.com/UmbrellaCrow612/fsearch/cli/out"
	"github.com/UmbrellaCrow612/fsearch/cli/shared"
	"github.com/UmbrellaCrow612/fsearch/cli/utils"
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

	if argsMap.Count {
		out.WriteToStdout(strconv.Itoa(len(matches)))
		out.ExitSuccess()
	}

	printMatchs(matches, argsMap)
	out.ExitSuccess()

	out.ExitSuccess()
}

const maxWorkers = 10

// Reads in a directory tree in parallel
func readInParallel(root string, argsMap *args.ArgsMap, searchTermRegex *regexp.Regexp) ([]shared.MatchEntry, error) {
	var matchEntrys []shared.MatchEntry
	var mu sync.Mutex
	var wg sync.WaitGroup
	var errs []error
	sem := make(chan struct{}, maxWorkers)

	modifiedBefore := utils.GetTimeValue(argsMap.ModifiedBefore)
	modifiedAfter := utils.GetTimeValue(argsMap.ModifiedAfter)

	sizeMultiplier := utils.GetSizeMultipler(argsMap)

	minSizeBytes := argsMap.MinSize * sizeMultiplier
	maxSizeBytes := argsMap.MaxSize * sizeMultiplier

	rootDepth := len(strings.Split(filepath.Clean(root), string(os.PathSeparator)))

	var reachedLimit atomic.Bool

	var read func(string, int, bool)
	read = func(path string, depth int, dirsOnly bool) {
		if argsMap.Limit > 0 && reachedLimit.Load() {
			return
		}

		sem <- struct{}{}
		defer func() { <-sem }()
		defer wg.Done()

		if argsMap.Depth > 0 && depth-rootDepth > argsMap.Depth {
			return
		}

		list, err := os.ReadDir(path)
		if err != nil {
			mu.Lock()
			errs = append(errs, fmt.Errorf("error reading %s: %w", path, err))
			mu.Unlock()
			return
		}

		for _, entry := range list {
			if argsMap.Limit > 0 && reachedLimit.Load() {
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
						if argsMap.Limit > 0 && len(matchEntrys) >= argsMap.Limit {
							mu.Unlock()
							reachedLimit.Store(true)
							return
						}
						matchEntrys = append(matchEntrys, shared.MatchEntry{Path: fullPath, Entry: entry})
						mu.Unlock()
					}
				}

				if argsMap.Depth <= 0 || depth-rootDepth < argsMap.Depth {
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
				mu.Lock()
				errs = append(errs, fmt.Errorf("error getting info for %s: %w", fullPath, err))
				mu.Unlock()
				continue
			}

			modTime := info.ModTime()
			size := info.Size()

			// The file extension ignoring the .
			ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(entry.Name())), ".")

			// --- Extension filters ---
			if len(argsMap.ExcludeExt) > 0 && slices.Contains(argsMap.ExcludeExt, ext) {
				continue
			}
			if len(argsMap.Ext) > 0 && !slices.Contains(argsMap.Ext, ext) {
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
			if argsMap.Limit > 0 && len(matchEntrys) >= argsMap.Limit {
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
