package args

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/UmbrellaCrow612/fsearch/src/out"
)

// Parses the args
func Parse() *ArgsMap {
	argMap := &ArgsMap{
		Type:           "file",
		Lines:          10,
		Depth:          0,
		Limit:          0,
		Preview:        false,
		Partial:        false,
		IgnoreCase:     false,
		Open:           false,
		Ext:            []string{},
		ExcludeExt:     []string{},
		ExcludeDir:     []string{},
		MinSize:        0,
		MaxSize:        0,
		ModifiedBefore: "Empty",
		ModifiedAfter:  "Empty",
		Hidden:         false,
		Count:          false,
		Stats:          false,
		Regex:          false,
		Debug:          false,
		Path:           "./",
		Term:           "test",
		SizeType:       "B",
	}
	setArgsMapValues(argMap)

	validateArgsMap(argMap)

	if argMap.Debug {
		printArgsMapValues(argMap)
		out.ExitSuccess()
	}

	return argMap
}

// Sets the cli args flags into the args map
func setArgsMapValues(argsMap *ArgsMap) {
	args := os.Args[1:]
	if len(args) < 2 {
		out.ExitError("Term or path not passed [..options..] [path]")
	}

	argsMap.Term = args[0]
	if isEmptyOrWhitespace(argsMap.Term) {
		out.ExitError("Search term cannot be empty or whitespace")
	}

	argsMap.Path = args[len(args)-1]
	if isEmptyOrWhitespace(argsMap.Path) {
		out.ExitError("Path cannot be empty or whitespace")
	}

	flags := args[1 : len(args)-1]

	for _, arg := range flags {
		switch {
		case arg == "--partial":
			argsMap.Partial = true
		case arg == "--ignore-case":
			argsMap.IgnoreCase = true
		case arg == "--open":
			argsMap.Open = true
		case arg == "--preview":
			argsMap.Preview = true

		case strings.HasPrefix(arg, "--lines="):
			val := strings.TrimPrefix(arg, "--lines=")
			lines, err := strconv.Atoi(val)
			if err != nil || lines < 0 {
				out.ExitError(fmt.Sprintf("Invalid value for --lines: '%s' (must be a positive integer)", val))
			}
			argsMap.Lines = lines

		case strings.HasPrefix(arg, "--limit="):
			val := strings.TrimPrefix(arg, "--limit=")
			limit, err := strconv.Atoi(val)
			if err != nil || limit <= 0 {
				out.ExitError(fmt.Sprintf("Invalid value for --limit: '%s' (must be a positive integer)", val))
			}
			argsMap.Limit = limit

		case strings.HasPrefix(arg, "--depth="):
			val := strings.TrimPrefix(arg, "--depth=")
			depth, err := strconv.Atoi(val)
			if err != nil || depth < 0 {
				out.ExitError(fmt.Sprintf("Invalid value for --depth: '%s' (must be a non-negative integer)", val))
			}
			argsMap.Depth = depth

		case strings.HasPrefix(arg, "--ext="):
			argsMap.Ext = strings.Split(strings.TrimPrefix(arg, "--ext="), ",")

		case strings.HasPrefix(arg, "--exclude-ext="):
			argsMap.ExcludeExt = strings.Split(strings.TrimPrefix(arg, "--exclude-ext="), ",")

		case strings.HasPrefix(arg, "--exclude-dir="):
			argsMap.ExcludeDir = strings.Split(strings.TrimPrefix(arg, "--exclude-dir="), ",")

		case strings.HasPrefix(arg, "--min-size="):
			val := strings.TrimPrefix(arg, "--min-size=")
			size, err := strconv.ParseInt(val, 10, 64)
			if err != nil || size < 0 {
				out.ExitError(fmt.Sprintf("Invalid value for --min-size: '%s' (must be a non-negative integer)", val))
			}
			argsMap.MinSize = size

		case strings.HasPrefix(arg, "--max-size="):
			val := strings.TrimPrefix(arg, "--max-size=")
			size, err := strconv.ParseInt(val, 10, 64)
			if err != nil || size < 0 {
				out.ExitError(fmt.Sprintf("Invalid value for --max-size: '%s' (must be a non-negative integer)", val))
			}
			argsMap.MaxSize = size

		case strings.HasPrefix(arg, "--size-type="):
			val := strings.ToUpper(strings.TrimPrefix(arg, "--size-type="))
			argsMap.SizeType = val

		case strings.HasPrefix(arg, "--modified-before="):
			val := strings.TrimPrefix(arg, "--modified-before=")
			if !isValidDate(val) {
				out.ExitError(fmt.Sprintf("Invalid date format for --modified-before: '%s' (expected YYYY-MM-DD)", val))
			}
			argsMap.ModifiedBefore = val

		case strings.HasPrefix(arg, "--modified-after="):
			val := strings.TrimPrefix(arg, "--modified-after=")
			if !isValidDate(val) {
				out.ExitError(fmt.Sprintf("Invalid date format for --modified-after: '%s' (expected YYYY-MM-DD)", val))
			}
			argsMap.ModifiedAfter = val

		case arg == "--hidden":
			argsMap.Hidden = true
		case arg == "--count":
			argsMap.Count = true
		case arg == "--stats":
			argsMap.Stats = true
		case arg == "--regex":
			argsMap.Regex = true
		case arg == "--debug":
			argsMap.Debug = true

		case strings.HasPrefix(arg, "--type="):
			val := strings.TrimPrefix(arg, "--type=")
			switch strings.ToLower(val) {
			case "file", "folder":
				argsMap.Type = val
			default:
				out.ExitError(fmt.Sprintf("Invalid value for --type: '%s' (expected 'file' or 'folder')", val))
			}

		default:
			out.ExitError("Unknown flag: " + arg)
		}
	}
}
