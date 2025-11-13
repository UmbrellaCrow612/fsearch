package utils

import "strings"

// IsHiddenFolderName checks if a folder name represents a hidden folder.
func IsHiddenFolderName(name string) bool {
	name = strings.TrimSpace(name)
	return len(name) > 1 && strings.HasPrefix(name, ".")
}