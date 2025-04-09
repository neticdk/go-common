package slice

import "github.com/neticdk/go-stdlib/xslices"

// FindFunc returns the first element in the slice that satisfies the
// predicate.
//
// It returns the default value for the type and false if no element
// satisfies the predicate.
// Deprecated: FindFunc is deprecated - use github.com/neticdk/go-stdlib/xslices.FindFunc
func FindFunc[T any](data []T, f func(T) bool) (T, bool) {
	return xslices.FindFunc(data, f)
}

// FindIFunc returns the index of the first element in the slice that satisfies
// the predicate.
//
// It returns -1 and false if no element satisfies the predicate
// Deprecated: FindIFunc is deprecated - use github.com/neticdk/go-stdlib/xslices.FindIFunc
func FindIFunc[T any](data []T, f func(T) bool) (int, bool) {
	return xslices.FindIFunc(data, f)
}
