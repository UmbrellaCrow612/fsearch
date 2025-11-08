package args

import (
	"fmt"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"

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

// IsValid checks if the OpenWith value is valid.
func (o OpenWith) IsValid() bool {
	return slices.Contains(AllOpenWith, o)
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

// IsValid checks if the SizeFormat value is valid.
func (s SizeFormat) IsValid() bool {
	return slices.Contains(AllSizeFormats, s)
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

// IsValid checks if the MatchType value is valid.
func (m MatchType) IsValid() bool {
	return slices.Contains(AllMatchTypes, m)
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
	argMap := &ArgsMap{}
	setArgsMapValues(argMap)

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
		out.ExitError("Search term cannot be empty or whitesapce")
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
			argsMap.Lines, _ = strconv.Atoi(strings.TrimPrefix(arg, "--lines="))
		case strings.HasPrefix(arg, "--limit="):
			argsMap.Limit, _ = strconv.Atoi(strings.TrimPrefix(arg, "--limit="))
		case strings.HasPrefix(arg, "--depth="):
			argsMap.Depth, _ = strconv.Atoi(strings.TrimPrefix(arg, "--depth="))
		case strings.HasPrefix(arg, "--ext="):
			argsMap.Ext = strings.Split(strings.TrimPrefix(arg, "--ext="), ",")
		case strings.HasPrefix(arg, "--exclude-ext="):
			argsMap.ExcludeExt = strings.Split(strings.TrimPrefix(arg, "--exclude-ext="), ",")
		case strings.HasPrefix(arg, "--exclude-dir="):
			argsMap.ExcludeDir = strings.Split(strings.TrimPrefix(arg, "--exclude-dir="), ",")
		case strings.HasPrefix(arg, "--min-size="):
			val := strings.TrimPrefix(arg, "--min-size=")
			argsMap.MinSize, argsMap.MinSizeFormat = parseSize(val)
		case strings.HasPrefix(arg, "--max-size="):
			val := strings.TrimPrefix(arg, "--max-size=")
			argsMap.MaxSize, argsMap.MaxSizeFormat = parseSize(val)
		case strings.HasPrefix(arg, "--modified-before="):
			argsMap.ModifiedBefore = strings.TrimPrefix(arg, "--modified-before=")
		case strings.HasPrefix(arg, "--modified-after="):
			argsMap.ModifiedAfter = strings.TrimPrefix(arg, "--modified-after=")
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
			argsMap.Type = MatchType(strings.TrimPrefix(arg, "--type="))
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
