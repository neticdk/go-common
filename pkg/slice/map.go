package slice

import "github.com/neticdk/go-stdlib/xslices"

// Map applies a function to each element of a slice and returns a new slice
// with the results.
// Deprecated: Map is deprecated - use github.com/neticdk/go-stdlib/xslices.Map
func Map[T, U any](data []T, f func(T) U) []U {
	return xslices.Map(data, f)
}
