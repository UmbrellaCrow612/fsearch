package shared

import "os"

// Represents a match either file or folder and it's info
type MatchEntry struct {
	Path  string
	Entry os.DirEntry
}
