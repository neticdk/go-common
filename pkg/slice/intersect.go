package slice

import "github.com/neticdk/go-stdlib/xslices"

// Intersect returns the intersection of two comparable slices.
// Deprecated: Intersect is deprecated - use github.com/neticdk/go-stdlib/xslices.Intersect
func Intersect[T comparable](a, b []T) []T {
	return xslices.Intersect(a, b)
}
