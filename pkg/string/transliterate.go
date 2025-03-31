package string

import "github.com/neticdk/go-stdlib/xstrings/transliterate"

// Transliterate converts a string to its transliterated form.
// Deprecated: Transliterate is deprecated and has been moved to github.com/neticdk/go-stdlib/xstrings/transliterate.String
func Transliterate(s string) string {
	return transliterate.String(s)
}
