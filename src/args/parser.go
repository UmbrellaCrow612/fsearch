package args

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/UmbrellaCrow612/fsearch/src/out"
)

// OpenWith defines supported applications that can open a file or folder.
type OpenWith string

// Supported enum-like values for OpenWith.
const (
	NotePad  OpenWith = "notepad.exe"
	VSCode   OpenWith = "code.exe"
	Explorer OpenWith = "explorer.exe"
)

// AllOpenWith contains all valid OpenWith values.
var AllOpenWith = []OpenWith{
	NotePad,
	VSCode,
	Explorer,
}

// String returns the string representation of the OpenWith value.
func (o OpenWith) String() string {
	return string(o)
}

// SizeFormat defines supported size units.
type SizeFormat string

// Supported enum-like values for SizeFormat.
const (
	Bytes SizeFormat = "B"
	KB    SizeFormat = "KB"
	MB    SizeFormat = "MB"
	GB    SizeFormat = "GB"
	TB    SizeFormat = "TB"
)

// AllSizeFormats contains all valid size formats.
var AllSizeFormats = []SizeFormat{
	Bytes,
	KB,
	MB,
	GB,
	TB,
}

// String returns the string representation of the SizeFormat value.
func (s SizeFormat) String() string {
	return string(s)
}

// MatchType defines whether to match files or folders.
type MatchType string

// Supported enum-like values for MatchType.
const (
	File   MatchType = "file"
	Folder MatchType = "folder"
)

// AllMatchTypes contains all valid MatchType values.
var AllMatchTypes = []MatchType{
	File,
	Folder,
}

// String returns the string representation of the MatchType value.
func (m MatchType) String() string {
	return string(m)
}

// Contains all the flag values mapped to their respective fields.
type ArgsMap struct {
	// If the term should be matched partially to the whole string as valid.
	Partial bool

	// If the term should be case insensitive matched.
	IgnoreCase bool

	// If the first match should be opened in a file explorer or text editor.
	Open bool

	// The exe to open the folder or file with.
	OpenWith OpenWith

	// If it should show a preview.
	Preview bool

	// Amount of lines to show from a file preview.
	Lines int

	// Number of items it will match before it will stop.
	Limit int

	// How many folders deep it will go from the root path passed.
	Depth int

	// List of extensions to match for.
	Ext []string

	// List of extensions to ignore.
	ExcludeExt []string

	// List of folders to ignore.
	ExcludeDir []string

	// Min size the match needs to be, just the number passed.
	MinSize int

	// The min size format such as MB, GB, etc.
	MinSizeFormat SizeFormat

	// Max size the match can be, just the number.
	MaxSize int

	// The format for the max size (e.g., KB, MB, GB, etc.).
	MaxSizeFormat SizeFormat

	// When the match can be modified before (date format: YYYY-MM-DD).
	ModifiedBefore string

	// When the match can be modified after (date format: YYYY-MM-DD).
	ModifiedAfter string

	// If it should match hidden files or folders (those starting with `.`).
	Hidden bool

	// If it should just output the count of matches.
	Count bool

	// If it should print a stats list.
	Stats bool

	// If the term should be treated as a regex pattern.
	Regex bool

	// If it should print debug information.
	Debug bool

	// If it should match type "file" or "folder".
	Type MatchType

	// The directory to scan
	Path string

	// The search term
	Term string
}

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
		OpenWith:       "notepad.exe",
		Ext:            []string{},
		ExcludeExt:     []string{},
		ExcludeDir:     []string{},
		MinSize:        0,
		MaxSize:        0,
		MinSizeFormat:  "B",
		MaxSizeFormat:  "B",
		ModifiedBefore: time.Now().Format("2006-01-02"),
		ModifiedAfter:  time.Now().Format("2006-01-02"),
		Hidden:         false,
		Count:          false,
		Stats:          false,
		Regex:          false,
		Debug:          false,
		Path:           "./",
		Term:           "test",
	}
	setArgsMapValues(argMap)

	err := validateArgsMap(argMap)
	if err != nil {
		out.ExitError(err.Error())
	}

	if argMap.Debug {
		printArgsMapValues(argMap)
		out.ExitSuccess()
	}

	return argMap
}

// validateArgsMap checks the parsed args for invalid combinations or values.
func validateArgsMap(argsMap *ArgsMap) error {
	absPath, err := filepath.Abs(argsMap.Path)
	if err != nil {
		return fmt.Errorf("could not resolve absolute path for %s: %v", argsMap.Path, err)
	}
	argsMap.Path = absPath

	fileInfo, err := os.Stat(argsMap.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("path does not exist: %s", argsMap.Path)
		}
		return fmt.Errorf("error checking path: %v", err)
	}
	if !fileInfo.IsDir() {
		return fmt.Errorf("path is not a directory: %s", argsMap.Path)
	}

	if argsMap.OpenWith != "" && !isValidEnumValue(string(argsMap.OpenWith), allOpenWithStrings()) {
		return fmt.Errorf("invalid --open-with value: %s. Must be one of: %v", argsMap.OpenWith, allOpenWithStrings())
	}
	if argsMap.Type != "" && !isValidEnumValue(string(argsMap.Type), allMatchTypeStrings()) {
		return fmt.Errorf("invalid --type value: %s. Must be one of: %v", argsMap.Type, allMatchTypeStrings())
	}

	if argsMap.Lines < 0 {
		return errors.New("--lines value cannot be negative")
	}
	if argsMap.Limit < 0 {
		return errors.New("--limit value cannot be negative")
	}
	if argsMap.Depth < 0 {
		return errors.New("--depth value cannot be negative")
	}
	if argsMap.MinSize < 0 {
		return errors.New("--min-size value cannot be negative")
	}
	if argsMap.MaxSize < 0 {
		return errors.New("--max-size value cannot be negative")
	}

	if argsMap.ModifiedBefore != "" && !isValidDate(argsMap.ModifiedBefore) {
		return fmt.Errorf("invalid date format for --modified-before: %s. Must be YYYY-MM-DD", argsMap.ModifiedBefore)
	}
	if argsMap.ModifiedAfter != "" && !isValidDate(argsMap.ModifiedAfter) {
		return fmt.Errorf("invalid date format for --modified-after: %s. Must be YYYY-MM-DD", argsMap.ModifiedAfter)
	}

	if argsMap.Regex {
		if _, err := regexp.Compile(argsMap.Term); err != nil {
			return fmt.Errorf("invalid regex term: %v", err)
		}
	}

	argsMap.Ext = cleanStringSlice(argsMap.Ext)
	argsMap.ExcludeExt = cleanStringSlice(argsMap.ExcludeExt)
	argsMap.ExcludeDir = cleanStringSlice(argsMap.ExcludeDir)

	return nil
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
		case strings.HasPrefix(arg, "--open-with="):
			argsMap.OpenWith = OpenWith(strings.TrimPrefix(arg, "--open-with="))
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
			size, format := parseSize(val)
			if size < 0 {
				out.ExitError(fmt.Sprintf("Invalid value for --min-size: '%s'", val))
			}
			argsMap.MinSize, argsMap.MinSizeFormat = size, format

		case strings.HasPrefix(arg, "--max-size="):
			val := strings.TrimPrefix(arg, "--max-size=")
			size, format := parseSize(val)
			if size < 0 {
				out.ExitError(fmt.Sprintf("Invalid value for --max-size: '%s'", val))
			}
			argsMap.MaxSize, argsMap.MaxSizeFormat = size, format

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
				argsMap.Type = MatchType(val)
			default:
				out.ExitError(fmt.Sprintf("Invalid value for --type: '%s' (expected 'file' or 'folder')", val))
			}

		default:
			out.ExitError("Unknown flag: " + arg)
		}
	}
}

// checks if a string is empty
func isEmptyOrWhitespace(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func parseSize(value string) (int, SizeFormat) {
	for _, unit := range AllSizeFormats {
		if strings.HasSuffix(value, unit.String()) {
			numStr := strings.TrimSuffix(value, unit.String())
			num, _ := strconv.Atoi(numStr)
			return num, unit
		}
	}
	return 0, Bytes
}

func printArgsMapValues(args *ArgsMap) {
	// Dereference the pointer
	v := reflect.ValueOf(args).Elem()
	t := v.Type()

	fmt.Println("----- ArgsMap Values -----")

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		if value.Kind() == reflect.Slice {
			fmt.Printf("%-15s | %-12s | %v\n", field.Name, value.Type(), sliceToString(value))
			continue
		}

		fmt.Printf("%-15s | %-12s | %v\n", field.Name, value.Type(), value.Interface())
	}

	fmt.Println("---------------------------")
}

func sliceToString(v reflect.Value) string {
	if v.Len() == 0 {
		return "[]"
	}
	s := "["
	for i := 0; i < v.Len(); i++ {
		s += fmt.Sprintf("%v", v.Index(i))
		if i < v.Len()-1 {
			s += ", "
		}
	}
	s += "]"
	return s
}

// --- Validation Helper Functions ---

// isValidEnumValue checks if a value exists in a list of valid strings.
func isValidEnumValue(value string, validValues []string) bool {
	return slices.Contains(validValues, value)
}

// allOpenWithStrings converts AllOpenWith to a string slice for validation.
func allOpenWithStrings() []string {
	s := make([]string, len(AllOpenWith))
	for i, v := range AllOpenWith {
		s[i] = v.String()
	}
	return s
}

// allMatchTypeStrings converts AllMatchTypes to a string slice for validation.
func allMatchTypeStrings() []string {
	s := make([]string, len(AllMatchTypes))
	for i, v := range AllMatchTypes {
		s[i] = v.String()
	}
	return s
}

// isValidDate checks if a string is a valid YYYY-MM-DD date.
func isValidDate(dateStr string) bool {
	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}

// cleanStringSlice removes empty strings from a slice,
// which can result from parsing "val1,,val2".
func cleanStringSlice(s []string) []string {
	if s == nil {
		return nil
	}
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	if len(r) == 0 {
		return nil // Return nil instead of an empty slice
	}
	return r
}
