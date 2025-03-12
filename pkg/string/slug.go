package string

import (
	"regexp"
	"strings"
)

// Slugify converts a string to a slug
func Slugify(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

	// Replace multiple spaces or hyphens with single hyphens.
	reSpaces := regexp.MustCompile(`[\s-]+`)
	s = reSpaces.ReplaceAllString(s, "-")

	// Replace characters other than alphanumeric and hyphen with an empty
	// string.
	reChars := regexp.MustCompile(`[^a-z0-9-]`)
	s = reChars.ReplaceAllString(s, "")

	// Remove leading and trailing hyphens.
	s = strings.Trim(s, "-")

	return s
}
