package string

import "github.com/neticdk/go-stdlib/xstrings"

// Slugify converts a string to a slug.
// By default, it converts the string to lowercase, transliterates it, and
// removes camel case.
// Deprecated: Slugify is deprecated and has been moved to github.com/neticdk/go-stdlib/xstrings.Slugify
func Slugify(s string, options ...xstrings.TransformOption) string {
	return xstrings.Slugify(s, options...)
}
