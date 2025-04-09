package slice

import "github.com/neticdk/go-stdlib/xslices"

// Filter returns a new slice containing only the elements of the slice that
// satisfy the predicate.
// Deprecated: Filter is deprecated - use github.com/neticdk/go-stdlib/xslices.Filter
func Filter[T any](data []T, f func(T) bool) []T {
	return xslices.Filter(data, f)
}
