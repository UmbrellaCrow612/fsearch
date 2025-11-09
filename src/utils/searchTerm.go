package utils

import (
	"fmt"
	"regexp"
	"strings"
)

// BuildSearchRegex builds a regex from a search term with options for partial and case-insensitive matching.
func BuildSearchRegex(term string, partial bool, caseInsensitive bool) (*regexp.Regexp, error) {
	if term == "" {
		return nil, fmt.Errorf("search term cannot be empty")
	}

	escaped := regexp.QuoteMeta(term)

	var patternBuilder strings.Builder
	if caseInsensitive {
		patternBuilder.WriteString("(?i)")
	}

	if partial {
		patternBuilder.WriteString(escaped)
	} else {
		patternBuilder.WriteString("^" + escaped + "$")
	}

	pattern := patternBuilder.String()

	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regex: %w", err)
	}

	return re, nil
}


// CompileRegex attempts to compile the given regex pattern.
// Returns the compiled *regexp.Regexp if successful, or an error if the pattern is invalid.
func CompileRegex(pattern string) (*regexp.Regexp, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern %q: %w", pattern, err)
	}
	return re, nil
}
