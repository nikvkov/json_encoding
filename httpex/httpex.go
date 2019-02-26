package httpex

import (
	"regexp"
)

// Does path match pattern?
func pathMatch(pattern, path string) bool {
	if len(pattern) == 0 {
		// should not happen
		return false
	}
	n := len(pattern)
	if pattern[n-1] != '/' {
		match, _ := regexp.MatchString(pattern, path)
		return match
	}
	fullMatch, _ := regexp.MatchString(pattern, string(path[0:n]))
	return len(path) >= n && fullMatch
}
