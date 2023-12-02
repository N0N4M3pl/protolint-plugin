package utils

import "strings"

// Contains reports whether one of substrs item is within s.
func ContainsStringFromArray(s string, substrs []string) bool {
	for _, substr := range substrs {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}
