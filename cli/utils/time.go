package utils

import "time"

// Converts a string of (YYYY-MM-DD) or Empty into a time object
// Retrusn either time or nil
func GetTimeValue(str string) *time.Time {
	if str != "" && str != "Empty" {
		t, err := time.Parse("2006-01-02", str)
		if err == nil {
			return &t
		}
	}

	return nil
}
