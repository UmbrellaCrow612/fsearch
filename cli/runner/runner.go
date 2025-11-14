package runner

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"

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

	matches, err := readDirectory(argsMap.Path, argsMap, searchTermRegex)
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
}

const maxWorkers = 10

// Reads in a directory tree
func readDirectory(root string, argsMap *args.ArgsMap, searchTermRegex *regexp.Regexp) ([]shared.MatchEntry, error) {
	var matchEntries []shared.MatchEntry

	modifiedBefore := utils.GetTimeValue(argsMap.ModifiedBefore)
	modifiedAfter := utils.GetTimeValue(argsMap.ModifiedAfter)
	sizeMultiplier := utils.GetSizeMultipler(argsMap)
	minSizeBytes := argsMap.MinSize * sizeMultiplier
	maxSizeBytes := argsMap.MaxSize * sizeMultiplier

	rootDepth := len(strings.Split(filepath.Clean(root), string(os.PathSeparator)))

	var walkFunc fs.WalkDirFunc
	walkFunc = func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if argsMap.Limit > 0 && len(matchEntries) >= argsMap.Limit {
			return filepath.SkipDir
		}

		depth := len(strings.Split(filepath.Clean(path), string(os.PathSeparator)))
		if argsMap.Depth > 0 && depth-rootDepth > argsMap.Depth {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		name := d.Name()
		fullPath := path

		if d.IsDir() {
			if !argsMap.Hidden && utils.IsHiddenFolderName(name) {
				return filepath.SkipDir
			}
			if slices.Contains(argsMap.ExcludeDir, name) {
				return filepath.SkipDir
			}

			if argsMap.Type == "folder" && searchTermRegex.MatchString(name) {
				matchEntries = append(matchEntries, shared.MatchEntry{Path: fullPath, Entry: d})
			}
			return nil
		}

		if argsMap.Type == "folder" {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		modTime := info.ModTime()
		size := info.Size()
		ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(name)), ".")

		if len(argsMap.ExcludeExt) > 0 && slices.Contains(argsMap.ExcludeExt, ext) {
			return nil
		}
		if len(argsMap.Ext) > 0 && !slices.Contains(argsMap.Ext, ext) {
			return nil
		}
		if modifiedBefore != nil && modTime.After(*modifiedBefore) {
			return nil
		}
		if modifiedAfter != nil && modTime.Before(*modifiedAfter) {
			return nil
		}
		if argsMap.MinSize > 0 && size < minSizeBytes {
			return nil
		}
		if argsMap.MaxSize > 0 && size > maxSizeBytes {
			return nil
		}
		if !searchTermRegex.MatchString(fullPath) {
			return nil
		}

		matchEntries = append(matchEntries, shared.MatchEntry{Path: fullPath, Entry: d})

		return nil
	}

	err := filepath.WalkDir(root, walkFunc)
	if err != nil {
		return matchEntries, err
	}

	return matchEntries, nil
}
