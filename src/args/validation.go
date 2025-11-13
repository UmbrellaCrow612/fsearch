package args

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/UmbrellaCrow612/fsearch/src/out"
)

// Validate args map values
func validateArgsMap(argsMap *ArgsMap) {
	// --- Validate path ---
	if isEmptyOrWhitespace(argsMap.Path) {
		out.ExitError("Path cannot be empty or whitespace")
	}

	_, err := os.Stat(argsMap.Path)
	if os.IsNotExist(err) {
		out.ExitError(fmt.Sprintf("The specified path does not exist: '%s'", argsMap.Path))
	}
	if err != nil {
		out.ExitError(fmt.Sprintf("Error accessing path '%s': %v", argsMap.Path, err))
	}

	absPath, err := filepath.Abs(argsMap.Path)
	if err != nil {
		out.ExitError(fmt.Sprintf("Failed to resolve absolute path for '%s': %v", argsMap.Path, err))
	}
	argsMap.Path = absPath
	// --- End validate path ---

	// --- Validate open-with ---
	if isEmptyOrWhitespace(argsMap.OpenWith) {
		out.ExitError("Open with cannot be empty or whitespace")
	}

	validViewers := getValidViewersForOS()

	found := false
	for _, viewer := range validViewers {
		if strings.EqualFold(viewer, argsMap.OpenWith) {
			found = true
			break
		}
	}

	if !found {
		if _, err := exec.LookPath(argsMap.OpenWith); err == nil {
			found = true
		}
	}

	if !found {
		out.ExitError(fmt.Sprintf(
			"Invalid program for --open-with: '%s'\nValid options for your OS include: %s",
			argsMap.OpenWith, strings.Join(validViewers, ", "),
		))
	}
	// --- End validate open-with ---

	// --- Validate lines ---
	if argsMap.Lines < 1 {
		out.ExitError("Lines cannot be below zero or zero")
	}
	// --- End Validate lines ---

	// --- Validate limit ---
	if argsMap.Limit < 0 {
		out.ExitError("Limit cannot be below zero")
	}
	// --- End Validate limit ---

	// --- Validate depth ---
	if argsMap.Depth < 0 {
		out.ExitError("Depth cannot be below zero")
	}
	// --- End Validate depth ---

	// --- Validation Ext ---
	filterEmptyStrings(&argsMap.Ext)
	// --- End Validation Ext ---

	// --- Validation ExcludeExt ---
	filterEmptyStrings(&argsMap.ExcludeExt)
	// --- End Validation ExcludeExt ---

	// --- Validation ExcludeDir ---
	filterEmptyStrings(&argsMap.ExcludeDir)
	// --- EndValidation ExcludeDir ---

	// --- Validation MinSize ---
	if argsMap.MinSize < 0 {
		out.ExitError("Min size cannot be below zero")
	}
	// --- End Validation MinSize ---

	// --- Validation MaxSize ---
	if argsMap.MaxSize < 0 {
		out.ExitError("Max size cannot be below zero")
	}
	// --- End Validation MaxSize ---

	// --- Validation SizeType ---
	if isEmptyOrWhitespace(argsMap.SizeType) {
		out.ExitError("Size typpe cannot be empty")
	}
	switch argsMap.SizeType {
	case "B", "KB", "MB", "GB":
	default:
		out.ExitError(fmt.Sprintf("Invalid value for --size-type: '%s' (expected B, KB, MB, or GB)", argsMap.SizeType))
	}
	// --- End Validation SizeType ---

	// --- Validation ModifiedBefore ---
	if !isValidDate(argsMap.ModifiedBefore) {
		out.ExitError("ModifiedBefore must be a valid DATE string (YYYY-MM-DD)")
	}
	// --- Validation ModifiedBefore End ---

	// --- Validation ModifiedAfter ---
	if !isValidDate(argsMap.ModifiedAfter) {
		out.ExitError("ModifiedAfter must be a valid DATE string (YYYY-MM-DD)")
	}
	// --- Validation ModifiedAfter End ---

	// --- Validation Type ---
	if !isValidType(argsMap.Type) {
		out.ExitError("Type must be either file or folder")
	}
	// --- Validation Type End---

	// --- Validation Term ---
	if isEmptyOrWhitespace(argsMap.Term) {
		out.ExitError("Term cannot be empty")
	}
	// --- Validation Term End---
}

// Checks if a string is a valid argmap type either "file" or "folder"
func isValidType(str string) bool {
	if str == "file" || str == "folder" {
		return true
	} else {
		return false
	}
}

// checks if a string is empty
func isEmptyOrWhitespace(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

// isValidDate checks if a string is a valid YYYY-MM-DD date.
func isValidDate(dateStr string) bool {
	if dateStr == "Empty" {
		return true
	}

	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}

// Removes empty or whitesapce
func filterEmptyStrings(input *[]string) {
	if input == nil {
		return
	}

	var filtered []string
	for _, s := range *input {
		if trimmed := strings.TrimSpace(s); trimmed != "" {
			filtered = append(filtered, trimmed)
		}
	}

	*input = filtered
}
