package args

// Contains all the flag values mapped to their respective fields.
type ArgsMap struct {
	// If the term should be matched partially to the whole string as valid.
	Partial bool

	// If the term should be case insensitive matched.
	IgnoreCase bool

	// If the first match should be opened in a file explorer or text editor.
	Open bool

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
	MinSize int64

	// Max size the match can be, just the number.
	MaxSize int64

	// A string of the type used for min max sizes either B KB MB or GB
	SizeType string

	// When the match can be modified before (date format: YYYY-MM-DD) or the string Empty for none values.
	ModifiedBefore string

	// When the match can be modified after (date format: YYYY-MM-DD) or the string Empty for none values.
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
	Type string

	// The directory to scan
	Path string

	// The search term
	Term string
}
